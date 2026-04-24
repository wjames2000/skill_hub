import * as vscode from 'vscode';
import { SkillTreeDataProvider } from './views/skillTreeView';
import { SkillHubService } from './services/skillHubService';
import { AuthService } from './services/authService';
import { AISkillPathDetector } from './services/aiSkillPathDetector';
import { registerCommands } from './commands';
import { checkRecommendations } from './views/recommendationHandler';

export async function activate(context: vscode.ExtensionContext) {
    console.log('SkillHub Pro extension is now active');

    const config = vscode.workspace.getConfiguration('skillhub');
    const apiEndpoint = config.get<string>('apiEndpoint') || 'https://api.skillhub.pro/v1';

    const authService = new AuthService(context);
    const skillHubService = new SkillHubService(authService, apiEndpoint);

    const pathDetector = new AISkillPathDetector();
    const detectedPaths = pathDetector.getDetectedPaths();
    if (detectedPaths.length > 0) {
        console.log(`Detected AI skill paths: ${detectedPaths.map(p => p.path).join(', ')}`);
    }

    const treeDataProvider = new SkillTreeDataProvider(skillHubService, authService);
    const treeView = vscode.window.createTreeView('skillhub-sidebar', {
        treeDataProvider,
        showCollapseAll: true
    });
    context.subscriptions.push(treeView);

    registerCommands(context, skillHubService, authService, treeDataProvider);

    const showRecommendations = config.get<boolean>('showRecommendations', true);
    if (showRecommendations) {
        context.subscriptions.push(
            vscode.window.onDidChangeActiveTextEditor(async (editor) => {
                if (editor && authService.isLoggedIn()) {
                    await checkRecommendations(editor, skillHubService);
                }
            })
        );
    }

    skillHubService.loadLocalSkills();
}

export function deactivate() {}