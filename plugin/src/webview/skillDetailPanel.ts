import * as vscode from 'vscode';
import { SkillHubService, Skill } from '../services/skillHubService';

export async function showSkillDetailPanel(
    context: vscode.ExtensionContext,
    skill: Skill,
    skillService: SkillHubService
): Promise<void> {
    const detail = await skillService.getSkillDetail(skill.id);
    if (!detail) {
        vscode.window.showErrorMessage('获取技能详情失败');
        return;
    }

    const html = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: var(--vscode-font-family); padding: 20px; }
        h1 { font-size: 1.5em; margin-bottom: 10px; }
        .meta { color: var(--vscode-descriptionForeground); font-size: 0.9em; margin-bottom: 20px; }
        .description { line-height: 1.6; margin-bottom: 20px; }
        .content { background: var(--vscode-editor-background); padding: 15px; border-radius: 4px; 
                 white-space: pre-wrap; font-family: var(--vscode-editor-font-family);
                 font-size: 0.9em; max-height: 400px; overflow: auto; }
        .btn { background: var(--vscode-button-background); color: var(--vscode-button-foreground);
               border: none; padding: 8px 16px; cursor: pointer; border-radius: 4px; margin-right: 10px; }
        .btn:hover { background: var(--vscode-button-hoverBackground); }
        .btn-secondary { background: transparent; border: 1px solid var(--vscode-button-secondaryBackground); }
    </style>
</head>
<body>
    <h1>${detail.name}</h1>
    <div class="meta">
        作者: ${detail.author || '未知'} | 
        ⭐ ${detail.github_stars} | 
        下载: ${detail.download_count} | 
        版本: ${detail.version}
    </div>
    <div class="description">${detail.description}</div>
    ${detail.content ? `<div class="content">${escapeHtml(detail.content.substring(0, 5000))}</div>` : ''}
    <div style="margin-top: 20px;">
        <button class="btn" onclick="installSkill()">安装技能</button>
        <button class="btn btn-secondary" onclick="viewOnline()">在SkillHub查看</button>
    </div>
    <script>
        const vscode = acquireVsCodeApi();
        function installSkill() {
            vscode.postMessage({ command: 'install', skill: ${JSON.stringify(skill)} });
        }
        function viewOnline() {
            vscode.postMessage({ command: 'openUrl', url: '${detail.source_url || ''}' });
        }
        function escapeHtml(text) {
            return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
        }
    </script>
</body>
</html>`;

    const panel = vscode.window.createWebviewPanel(
        'skillhub.detail',
        detail.name,
        vscode.ViewColumn.One,
        { enableScripts: true }
    );

    panel.webview.html = html;

    panel.webview.onDidReceiveMessage(async (message) => {
        if (message.command === 'install') {
            await skillService.installSkill(message.skill);
            panel.dispose();
        } else if (message.command === 'openUrl') {
            vscode.env.openExternal(vscode.Uri.parse(message.url));
        }
    });
}

function escapeHtml(text: string): string {
    return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}