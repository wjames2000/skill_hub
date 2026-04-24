import * as vscode from 'vscode';
import { SkillHubService, Skill, LocalSkill } from '../services/skillHubService';
import { AuthService } from '../services/authService';
import { SkillTreeDataProvider } from '../views/skillTreeView';
import { showSkillDetailPanel } from '../webview/skillDetailPanel';
import { showRecommendationNotification } from '../views/recommendationHandler';

export function registerCommands(
    context: vscode.ExtensionContext,
    skillService: SkillHubService,
    authService: AuthService,
    treeDataProvider: SkillTreeDataProvider
): void {
    context.subscriptions.push(
        vscode.commands.registerCommand('skillhub.search', async () => {
            const query = await vscode.window.showInputBox({
                prompt: '搜索技能',
                placeHolder: '输入关键词...'
            });

            if (query) {
                const results = await skillService.searchSkills(query);
                if (results.length === 0) {
                    vscode.window.showInformationMessage('未找到相关技能');
                } else if (results.length === 1) {
                    await showSkillDetailPanel(context, results[0], skillService);
                } else {
                    const selected = await vscode.window.showQuickPick(
                        results.map(s => ({ label: s.name, description: s.description, skill: s })),
                        { placeHolder: '选择技能' }
                    );
                    if (selected?.skill) {
                        await showSkillDetailPanel(context, selected.skill, skillService);
                    }
                }
            }
        }),

        vscode.commands.registerCommand('skillhub.login', async () => {
            if (authService.isLoggedIn()) {
                const user = authService.getUser();
                vscode.window.showInformationMessage(`已登录: ${user?.name}`);
            } else {
                const result = await vscode.window.showQuickPick(
                    [
                        { label: 'GitHub OAuth 登录', method: 'github' },
                        { label: '使用 API Key', method: 'apikey' }
                    ],
                    { placeHolder: '选择登录方式' }
                );

                if (result?.method === 'github') {
                    await authService.loginWithGitHub();
                } else if (result?.method === 'apikey') {
                    const apiKey = await vscode.window.showInputBox({
                        prompt: '输入 API Key',
                        placeHolder: 'sk-...'
                    });
                    if (apiKey) {
                        // API Key 登录逻辑
                        vscode.window.showInformationMessage('API Key 登录功能开发中');
                    }
                }
            }
        }),

        vscode.commands.registerCommand('skillhub.logout', async () => {
            await authService.logout();
            vscode.window.showInformationMessage('已退出登录');
            treeDataProvider.refresh();
        }),

        vscode.commands.registerCommand('skillhub.refresh', () => {
            treeDataProvider.refresh();
        }),

        vscode.commands.registerCommand('skillhub.install', async (skill: Skill) => {
            await skillService.installSkill(skill);
        }),

        vscode.commands.registerCommand('skillhub.showDetail', async (skill: Skill) => {
            await showSkillDetailPanel(context, skill, skillService);
        }),

        vscode.commands.registerCommand('skillhub.openSkillFolder', async (skill: LocalSkill) => {
            skillService.openSkillFolder(skill);
        }),

        vscode.commands.registerCommand('skillhub.toggleSkill', async (skill: LocalSkill) => {
            await skillService.toggleSkill(skill);
        }),

        vscode.commands.registerCommand('skillhub.checkUpdates', async () => {
            const hasUpdate = await skillService.checkUpdates();
            if (hasUpdate.length === 0) {
                vscode.window.showInformationMessage('所有技能已是最新版本');
            } else {
                const selected = await vscode.window.showQuickPick(
                    hasUpdate.map(s => ({ label: s.name, description: `新版本: ${s.version}`, skill: s })),
                    { placeHolder: '选择要更新的技能' }
                );
                if (selected?.skill) {
                    await skillService.updateSkill(selected.skill);
                }
            }
        }),

        vscode.commands.registerCommand('skillhub.syncToCloud', async () => {
            if (!authService.isLoggedIn()) {
                vscode.window.showErrorMessage('请先登录');
                return;
            }

            const localSkills = skillService.getLocalSkills();
            if (localSkills.length === 0) {
                vscode.window.showInformationMessage('没有本地技能需要同步');
                return;
            }

            vscode.window.showInformationMessage(`正在同步 ${localSkills.length} 个技能...`);
            vscode.window.showInformationMessage('同步功能开发中');
        }),

        vscode.commands.registerCommand('skillhub.viewPopular', () => {
            treeDataProvider.setViewMode('popular');
        }),

        vscode.commands.registerCommand('skillhub.viewLocal', () => {
            treeDataProvider.setViewMode('local');
        }),

        vscode.commands.registerCommand('skillhub.viewCategory', () => {
            treeDataProvider.setViewMode('category');
        })
    );
}