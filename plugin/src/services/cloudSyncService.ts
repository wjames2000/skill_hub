import * as vscode from 'vscode';
import axios, { AxiosInstance } from 'axios';
import { AuthService } from './authService';
import { LocalSkill } from './skillHubService';

export interface CloudSkill {
    skill_id: string;
    local_version: string;
    enabled: boolean;
    installed_at: string;
    last_synced_at: string;
}

export interface SyncResult {
    uploaded: number;
    downloaded: number;
    conflicts: CloudSkill[];
}

export class CloudSyncService {
    private apiClient: AxiosInstance;

    constructor(
        private authService: AuthService,
        private apiEndpoint: string
    ) {
        this.apiClient = axios.create({
            baseURL: apiEndpoint,
            timeout: 10000
        });
    }

    async getCloudSkills(): Promise<CloudSkill[]> {
        try {
            const token = await this.authService.getToken();
            if (!token) throw new Error('Not logged in');

            const response = await this.apiClient.get('/plugin/user/installed', {
                headers: { Authorization: `Bearer ${token}` }
            });
            return response.data.data?.items || [];
        } catch (error) {
            console.error('Failed to get cloud skills:', error);
            return [];
        }
    }

    async syncToCloud(localSkills: LocalSkill[]): Promise<SyncResult> {
        const result: SyncResult = {
            uploaded: 0,
            downloaded: 0,
            conflicts: []
        };

        if (!this.authService.isLoggedIn()) {
            throw new Error('请先登录');
        }

        const cloudSkills = await this.getCloudSkills();
        const cloudMap = new Map(cloudSkills.map(s => [s.skill_id, s]));

        for (const local of localSkills) {
            const cloud = cloudMap.get(local.id);

            if (!cloud) {
                await this.uploadSkillState(local);
                result.uploaded++;
            } else if (cloud.local_version !== local.installedVersion) {
                result.conflicts.push(cloud);
            }
        }

        for (const cloud of cloudSkills) {
            const local = localSkills.find(s => s.id === cloud.skill_id);
            if (!local) {
                result.downloaded++;
            }
        }

        return result;
    }

    async uploadSkillState(skill: LocalSkill): Promise<void> {
        const token = await this.authService.getToken();
        if (!token) throw new Error('Not logged in');

        await this.apiClient.post('/plugin/user/installed', {
            skill_id: skill.id,
            action: 'install',
            local_version: skill.installedVersion,
            enabled: skill.enabled
        }, {
            headers: { Authorization: `Bearer ${token}` }
        });
    }

    async removeFromCloud(skillId: string): Promise<void> {
        const token = await this.authService.getToken();
        if (!token) throw new Error('Not logged in');

        await this.apiClient.post('/plugin/user/installed', {
            skill_id: skillId,
            action: 'uninstall'
        }, {
            headers: { Authorization: `Bearer ${token}` }
        });
    }

    async resolveConflict(skillId: string, useCloud: boolean): Promise<void> {
        const token = await this.authService.getToken();
        if (!token) throw new Error('Not logged in');

        if (useCloud) {
            await this.apiClient.post('/plugin/user/conflicts/resolve', {
                skill_id: skillId,
                resolution: 'cloud'
            }, {
                headers: { Authorization: `Bearer ${token}` }
            });
        } else {
            const cloudSkills = await this.getCloudSkills();
            const cloud = cloudSkills.find(s => s.skill_id === skillId);
            if (cloud) {
                await this.uploadSkillState({
                    id: skillId,
                    installedVersion: cloud.local_version,
                    enabled: cloud.enabled
                } as LocalSkill);
            }
        }
    }

    async lastSyncTime(): Promise<string | null> {
        const key = 'skillhub.sync.lastSyncTime';
        return vscode.workspace.getConfiguration('skillhub').get<string>(key) || null;
    }

    async updateSyncTime(): Promise<void> {
        const key = 'skillhub.sync.lastSyncTime';
        const config = vscode.workspace.getConfiguration('skillhub');
        await config.update(key, new Date().toISOString());
    }
}