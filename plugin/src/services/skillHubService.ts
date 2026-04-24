import * as vscode from 'vscode';
import axios, { AxiosInstance } from 'axios';
import { exec } from 'child_process';
import { promisify } from 'util';
import * as path from 'path';
import * as fs from 'fs/promises';
import { AuthService } from './authService';

const execAsync = promisify(exec);

export interface Skill {
    id: string;
    name: string;
    description: string;
    version: string;
    author: string;
    github_stars: number;
    download_count: number;
    source_url: string;
    skill_md_url?: string;
    category?: string;
}

export interface LocalSkill extends Skill {
    path: string;
    enabled: boolean;
    installedVersion: string;
    hasUpdate: boolean;
    lastUpdated?: string;
}

export interface Category {
    id: string;
    name: string;
    icon?: string;
    parent_id?: string;
}

export class SkillHubService {
    private apiClient: AxiosInstance;
    private localSkillsPath: string;
    private localSkills: Map<string, LocalSkill> = new Map();
    private onLocalSkillsChange?: () => void;

    constructor(private authService: AuthService, private apiEndpoint: string) {
        this.apiClient = axios.create({
            baseURL: apiEndpoint,
            timeout: 10000
        });

        this.apiClient.interceptors.request.use(async (config) => {
            const token = await this.authService.getToken();
            if (token) {
                config.headers.Authorization = `Bearer ${token}`;
            }
            return config;
        });

        const config = vscode.workspace.getConfiguration('skillhub');
        this.localSkillsPath = config.get<string>('localSkillsPath') || '.claude/skills';
    }

    setOnLocalSkillsChange(callback: () => void) {
        this.onLocalSkillsChange = callback;
    }

    async getPopularSkills(): Promise<Skill[]> {
        try {
            const response = await this.apiClient.get('/plugin/skills/popular');
            return response.data.data?.items || [];
        } catch (error) {
            console.error('Failed to get popular skills:', error);
            return [];
        }
    }

    async getCategories(): Promise<Category[]> {
        try {
            const response = await this.apiClient.get('/categories');
            return response.data.data?.items || [];
        } catch (error) {
            console.error('Failed to get categories:', error);
            return [];
        }
    }

    async getSkillsByCategory(categoryId: string): Promise<Skill[]> {
        try {
            const response = await this.apiClient.get('/skills', {
                params: { category_id: categoryId }
            });
            return response.data.data?.items || [];
        } catch (error) {
            console.error('Failed to get skills by category:', error);
            return [];
        }
    }

    async searchSkills(query: string): Promise<Skill[]> {
        try {
            const response = await this.apiClient.post('/skills/search', { query });
            return response.data.data?.items || [];
        } catch (error) {
            console.error('Failed to search skills:', error);
            return [];
        }
    }

    async getSkillDetail(id: string): Promise<(Skill & { content?: string }) | null> {
        try {
            const response = await this.apiClient.get(`/skills/${id}`);
            return response.data.data;
        } catch (error) {
            console.error('Failed to get skill detail:', error);
            return null;
        }
    }

    async installSkill(skill: Skill): Promise<boolean> {
        const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
        if (!workspaceRoot) {
            vscode.window.showErrorMessage('请先打开一个工作区');
            return false;
        }

        const targetDir = path.join(workspaceRoot, this.localSkillsPath, skill.id);

        return vscode.window.withProgress({
            location: vscode.ProgressLocation.Notification,
            title: `正在安装 ${skill.name}...`,
            cancellable: false
        }, async (progress) => {
            try {
                progress.report({ message: '获取下载信息...' });
                const resp = await this.apiClient.get(`/plugin/skills/${skill.id}/download`);
                const { download_url, type } = resp.data.data;

                progress.report({ message: '正在下载...' });
                await fs.mkdir(targetDir, { recursive: true });

                if (type === 'git') {
                    const repoUrl = download_url.replace(/\.git$/, '');
                    await execAsync(`git clone ${repoUrl} "${targetDir}"`, { cwd: path.dirname(targetDir) });
                } else {
                    const zipPath = path.join(targetDir, 'skill.zip');
                    const response = await axios.get(download_url, { responseType: 'arraybuffer' });
                    await fs.writeFile(zipPath, Buffer.from(response.data));
                    await execAsync(`unzip -o "${zipPath}" -d "${targetDir}"`);
                    await fs.unlink(zipPath);
                }

                this.saveLocalSkillMeta(skill, targetDir);

                if (this.authService.isLoggedIn()) {
                    try {
                        await this.apiClient.post('/plugin/user/installed', {
                            skill_id: skill.id,
                            action: 'install',
                            local_version: skill.version
                        });
                    } catch (e) {
                        console.error('Failed to sync to cloud:', e);
                    }
                }

                this.loadLocalSkills();
                vscode.window.showInformationMessage(`${skill.name} 安装成功！`);
                return true;
            } catch (error: any) {
                vscode.window.showErrorMessage(`安装失败: ${error.message}`);
                try {
                    await fs.rm(targetDir, { recursive: true, force: true });
                } catch {}
                return false;
            }
        });
    }

    async updateSkill(skill: LocalSkill): Promise<boolean> {
        const fullSkill: Skill = {
            id: skill.id,
            name: skill.name,
            description: skill.description,
            version: skill.version,
            author: skill.author,
            github_stars: skill.github_stars,
            download_count: skill.download_count,
            source_url: skill.source_url
        };
        return this.installSkill(fullSkill);
    }

    async uninstallSkill(skill: LocalSkill): Promise<boolean> {
        try {
            await fs.rm(skill.path, { recursive: true, force: true });
            this.localSkills.delete(skill.id);

            if (this.authService.isLoggedIn()) {
                try {
                    await this.apiClient.post('/plugin/user/installed', {
                        skill_id: skill.id,
                        action: 'uninstall'
                    });
                } catch (e) {
                    console.error('Failed to sync to cloud:', e);
                }
            }

            this.loadLocalSkills();
            vscode.window.showInformationMessage(`${skill.name} 已卸载`);
            return true;
        } catch (error: any) {
            vscode.window.showErrorMessage(`卸载失败: ${error.message}`);
            return false;
        }
    }

    async toggleSkill(skill: LocalSkill): Promise<void> {
        const newPath = skill.enabled
            ? skill.path + '.disabled'
            : skill.path.replace(/\.disabled$/, '');

        try {
            await fs.rename(skill.path, newPath);
            skill.enabled = !skill.enabled;
            skill.path = newPath;
            this.saveLocalSkillMeta(skill, newPath);
            this.loadLocalSkills();

            const msg = skill.enabled ? '启用' : '禁用';
            vscode.window.showInformationMessage(`${skill.name} 已${msg}`);
        } catch (error: any) {
            vscode.window.showErrorMessage(`操作失败: ${error.message}`);
        }
    }

    loadLocalSkills(): void {
        this.localSkills.clear();
        const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
        if (!workspaceRoot) return;

        const skillsDir = path.join(workspaceRoot, this.localSkillsPath);
        this.scanDirectory(skillsDir);

        const disabledDir = path.join(workspaceRoot, this.localSkillsPath + '.disabled');
        this.scanDirectory(disabledDir, false);

        if (this.onLocalSkillsChange) {
            this.onLocalSkillsChange();
        }
    }

    private async scanDirectory(dir: string, enabled: boolean = true): Promise<void> {
        try {
            const entries = await fs.readdir(dir, { withFileTypes: true });
            for (const entry of entries) {
                if (!entry.isDirectory()) continue;

                const metaPath = path.join(dir, entry.name, '.skillhub.json');
                try {
                    const content = await fs.readFile(metaPath, 'utf-8');
                    const meta = JSON.parse(content);
                    this.localSkills.set(meta.id, {
                        ...meta,
                        path: path.join(dir, entry.name),
                        enabled,
                        installedVersion: meta.version,
                        hasUpdate: false,
                        github_stars: 0,
                        download_count: 0,
                        source_url: ''
                    });
                } catch {
                    // No meta file, skip
                }
            }
        } catch {
            // Directory doesn't exist
        }
    }

    getLocalSkills(): LocalSkill[] {
        return Array.from(this.localSkills.values());
    }

    async checkUpdates(): Promise<LocalSkill[]> {
        const localSkills = this.getLocalSkills();
        const skillsWithUpdate: LocalSkill[] = [];

        for (const skill of localSkills) {
            try {
                const resp = await this.apiClient.get(`/skills/${skill.id}`);
                const remoteVersion = resp.data.data?.version;
                if (remoteVersion && remoteVersion !== skill.installedVersion) {
                    skill.hasUpdate = true;
                    skillsWithUpdate.push(skill);
                }
            } catch {
                // Skip on error
            }
        }

        return skillsWithUpdate;
    }

    async getRecommendations(filePath: string, language: string, cursorContext: string): Promise<Skill[]> {
        try {
            const response = await this.apiClient.post('/plugin/recommend', {
                file_path: filePath,
                file_language: language,
                cursor_context: cursorContext,
                top_k: 3
            });
            return response.data.data?.recommendations || [];
        } catch (error) {
            console.error('Failed to get recommendations:', error);
            return [];
        }
    }

    openSkillFolder(skill: LocalSkill): void {
        vscode.commands.executeCommand('revealFileInOS', vscode.Uri.file(skill.path));
    }

    private async saveLocalSkillMeta(skill: Skill, localPath: string): Promise<void> {
        const metaPath = path.join(localPath, '.skillhub.json');
        const meta = {
            id: skill.id,
            name: skill.name,
            version: skill.version,
            author: skill.author,
            description: skill.description,
            installedAt: new Date().toISOString()
        };
        await fs.writeFile(metaPath, JSON.stringify(meta, null, 2));
    }
}