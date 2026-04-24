import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs/promises';

export interface AISkillPath {
    name: string;
    path: string;
    type: 'claude' | 'openai' | 'custom';
    exists: boolean;
}

export class AISkillPathDetector {
    private workspaceRoot: string | undefined;
    private detectedPaths: AISkillPath[] = [];

    constructor() {
        this.workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
        if (this.workspaceRoot) {
            this.detectPaths();
        }
    }

    private async detectPaths(): Promise<void> {
        if (!this.workspaceRoot) return;

        const defaultPaths: AISkillPath[] = [
            { name: 'Claude Skills', path: '.claude/skills', type: 'claude', exists: false },
            { name: 'Claude Commands', path: '.claude/commands', type: 'claude', exists: false },
            { name: 'OpenAI GPTs', path: 'gpts', type: 'openai', exists: false },
            { name: 'Custom Skills', path: 'skills', type: 'custom', exists: false },
            { name: 'AI Plugins', path: '.ai/plugins', type: 'custom', exists: false },
        ];

        for (const p of defaultPaths) {
            const fullPath = path.join(this.workspaceRoot, p.path);
            try {
                await fs.access(fullPath);
                p.exists = true;
            } catch {
                p.exists = false;
            }
        }

        this.detectedPaths = defaultPaths.filter(p => p.exists);

        await this.scanForAdditionalPaths();
    }

    private async scanForAdditionalPaths(): Promise<void> {
        if (!this.workspaceRoot) return;

        const searchPatterns = [
            '.*-skills',
            '.*-commands', 
            'ai-*',
            '*assistant*',
        ];

        try {
            const entries = await fs.readdir(this.workspaceRoot);
            for (const entry of entries) {
                const fullPath = path.join(this.workspaceRoot!, entry);
                try {
                    const stat = await fs.stat(fullPath);
                    if (stat.isDirectory()) {
                        const skillFile = path.join(fullPath, 'SKILL.md');
                        const readmeFile = path.join(fullPath, 'README.md');
                        try {
                            await fs.access(skillFile);
                            this.detectedPaths.push({
                                name: `${entry} Skills`,
                                path: entry,
                                type: 'custom',
                                exists: true
                            });
                        } catch {
                            try {
                                await fs.access(readmeFile);
                                const content = await fs.readFile(readmeFile, 'utf-8');
                                if (content.toLowerCase().includes('skill')) {
                                    this.detectedPaths.push({
                                        name: `${entry} (detected)`,
                                        path: entry,
                                        type: 'custom',
                                        exists: true
                                    });
                                }
                            } catch {
                                // Neither found
                            }
                        }
                    }
                } catch {
                    // Not accessible
                }
            }
        } catch {
            // Workspace not accessible
        }
    }

    getDetectedPaths(): AISkillPath[] {
        return this.detectedPaths;
    }

    getPrimarySkillPath(): string | undefined {
        if (this.detectedPaths.length === 0) return undefined;

        const claudePath = this.detectedPaths.find(p => p.type === 'claude');
        if (claudePath) return claudePath.path;

        return this.detectedPaths[0].path;
    }

    getSkillPathsForType(type: 'claude' | 'openai' | 'custom'): AISkillPath[] {
        return this.detectedPaths.filter(p => p.type === type);
    }

    async registerNewPath(name: string, pathPattern: string): Promise<AISkillPath> {
        const fullPath = path.join(this.workspaceRoot || '', pathPattern);
        try {
            await fs.mkdir(fullPath, { recursive: true });
            const newPath: AISkillPath = {
                name,
                path: pathPattern,
                type: 'custom',
                exists: true
            };
            this.detectedPaths.push(newPath);
            return newPath;
        } catch (error) {
            throw new Error(`无法创建路径: ${error}`);
        }
    }
}