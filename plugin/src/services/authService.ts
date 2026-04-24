import * as vscode from 'vscode';

const TOKEN_KEY = 'skillhub.auth.token';
const USER_KEY = 'skillhub.auth.user';

export interface AuthUser {
    id: string;
    email: string;
    name: string;
    avatar_url?: string;
}

export class AuthService {
    private token: string | null = null;
    private user: AuthUser | null = null;

    constructor(private context: vscode.ExtensionContext) {
        this.loadFromStorage();
    }

    private loadFromStorage(): void {
        this.token = this.context.globalState.get<string>(TOKEN_KEY) || null;
        this.user = this.context.globalState.get<AuthUser>(USER_KEY) || null;
    }

    async getToken(): Promise<string | null> {
        return this.token;
    }

    isLoggedIn(): boolean {
        return !!this.token && !!this.user;
    }

    getUser(): AuthUser | null {
        return this.user;
    }

    async login(token: string, user: AuthUser): Promise<void> {
        this.token = token;
        this.user = user;
        await this.context.globalState.update(TOKEN_KEY, token);
        await this.context.globalState.update(USER_KEY, user);
    }

    async logout(): Promise<void> {
        this.token = null;
        this.user = null;
        await this.context.globalState.update(TOKEN_KEY, undefined);
        await this.context.globalState.update(USER_KEY, undefined);
    }

    async loginWithGitHub(): Promise<boolean> {
        const config = vscode.workspace.getConfiguration('skillhub');
        const apiEndpoint = config.get<string>('apiEndpoint') || 'https://api.skillhub.pro/v1';

        const oauthUrl = `${apiEndpoint.replace('/v1', '')}/auth/github/start`;
        const callbackUri = `vscode://skillhub/auth/callback`;

        try {
            await vscode.commands.executeCommand('vscode.open', vscode.Uri.parse(oauthUrl));

            const token = await vscode.window.showInputBox({
                prompt: '请输入GitHub OAuth验证代码',
                placeHolder: '粘贴验证URL中的code参数值'
            });

            if (!token) return false;

            const response = await fetch(`${apiEndpoint}/auth/github/callback`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ code: token })
            });

            if (!response.ok) {
                vscode.window.showErrorMessage('GitHub登录失败');
                return false;
            }

            const data = await response.json();
            if (data.data?.token && data.data?.user) {
                await this.login(data.data.token, data.data.user);
                vscode.window.showInformationMessage(`欢迎 ${data.data.user.name}!`);
                return true;
            }

            return false;
        } catch (error: any) {
            vscode.window.showErrorMessage(`登录失败: ${error.message}`);
            return false;
        }
    }
}