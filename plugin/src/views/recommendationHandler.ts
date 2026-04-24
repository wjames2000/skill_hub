import * as vscode from 'vscode';
import { SkillHubService, Skill } from '../services/skillHubService';

export async function checkRecommendations(
    editor: vscode.TextEditor,
    skillService: SkillHubService
): Promise<void> {
    const document = editor.document;
    const cursorPosition = editor.selection.active;
    const line = document.lineAt(cursorPosition.line);
    const cursorContext = line.text;

    if (!cursorContext.trim()) return;

    try {
        const recommendations = await skillService.getRecommendations(
            document.fileName,
            document.languageId,
            cursorContext
        );

        if (recommendations.length > 0) {
            await showRecommendationNotification(recommendations);
        }
    } catch (error) {
        console.error('Failed to get recommendations:', error);
    }
}

export async function showRecommendationNotification(recommendations: Skill[]): Promise<void> {
    if (recommendations.length === 0) return;

    const items = recommendations.map(s => ({
        label: s.name,
        description: s.description?.substring(0, 50) + '...',
        skill: s
    }));

    const selected = await vscode.window.showQuickPick(items, {
        placeHolder: '推荐技能：',
        ignoreFocusOut: true
    });

    if (selected?.skill) {
        vscode.commands.executeCommand('skillhub.showDetail', selected.skill);
    }
}