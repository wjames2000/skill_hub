import * as vscode from 'vscode';
import { SkillHubService, Skill, LocalSkill, Category } from '../services/skillHubService';
import { AuthService } from '../services/authService';

export class SkillTreeItem extends vscode.TreeItem {
    constructor(
        public readonly label: string,
        public readonly collapsibleState: vscode.TreeItemCollapsibleState,
        public readonly skill?: Skill,
        public readonly localSkill?: LocalSkill,
        public readonly category?: Category
    ) {
        super(label, collapsibleState);

        if (skill) {
            this.id = skill.id;
            this.description = `${skill.github_stars} ⭐`;
            this.iconPath = new vscode.ThemeIcon('symbol-property');
            this.contextValue = 'skill';
        } else if (localSkill) {
            this.id = localSkill.id;
            this.description = localSkill.enabled ? '已启用' : '已禁用';
            this.iconPath = new vscode.ThemeIcon(localSkill.enabled ? 'check' : 'circle-slash');
            this.contextValue = 'localSkill';
            if (localSkill.hasUpdate) {
                this.iconPath = new vscode.ThemeIcon('arrow-down');
                this.description = '有更新';
            }
        } else if (category) {
            this.id = category.id;
            this.iconPath = new vscode.ThemeIcon('symbol-method');
            this.contextValue = 'category';
        }
    }
}

export class SkillTreeDataProvider implements vscode.TreeDataProvider<SkillTreeItem> {
    private _onDidChangeTreeData = new vscode.EventEmitter<SkillTreeItem | undefined>();
    readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

    private popularSkills: Skill[] = [];
    private categories: Category[] = [];
    private categorySkills: Map<string, Skill[]> = new Map();
    private localSkills: LocalSkill[] = [];
    private viewMode: 'popular' | 'local' | 'category' = 'popular';

    constructor(
        private skillService: SkillHubService,
        private authService: AuthService
    ) {
        this.loadData();
        skillService.setOnLocalSkillsChange(() => this.refresh());
    }

    async loadData(): Promise<void> {
        this.localSkills = this.skillService.getLocalSkills();

        if (this.viewMode === 'popular') {
            this.popularSkills = await this.skillService.getPopularSkills();
        } else if (this.viewMode === 'category') {
            this.categories = await this.skillService.getCategories();
        }

        this._onDidChangeTreeData.fire(undefined);
    }

    refresh(): void {
        this.loadData();
    }

    setViewMode(mode: 'popular' | 'local' | 'category'): void {
        this.viewMode = mode;
        this.loadData();
    }

    getTreeItem(element: SkillTreeItem): SkillTreeItem {
        return element;
    }

    async getChildren(element?: SkillTreeItem): Promise<SkillTreeItem[]> {
        if (!element) {
            return this.getRootItems();
        }

        if (element.category) {
            const skills = await this.skillService.getSkillsByCategory(element.category.id);
            return skills.map(s => new SkillTreeItem(s.name, vscode.TreeItemCollapsibleState.None, s));
        }

        if (element.label === '本地技能') {
            return this.localSkills.map(s => new SkillTreeItem(
                s.name,
                vscode.TreeItemCollapsibleState.None,
                undefined,
                s
            ));
        }

        if (element.label === '热门技能') {
            return this.popularSkills.map(s => new SkillTreeItem(
                s.name,
                vscode.TreeItemCollapsibleState.None,
                s
            ));
        }

        if (element.label === '分类浏览') {
            return this.categories.map(c => new SkillTreeItem(
                c.name,
                vscode.TreeItemCollapsibleState.Collapsed,
                undefined,
                undefined,
                c
            ));
        }

        return [];
    }

    private async getRootItems(): Promise<SkillTreeItem[]> {
        const items: SkillTreeItem[] = [];

        items.push(new SkillTreeItem('热门技能', vscode.TreeItemCollapsibleState.Collapsed));

        items.push(new SkillTreeItem('分类浏览', vscode.TreeItemCollapsibleState.Collapsed));

        if (this.localSkills.length > 0 || this.authService.isLoggedIn()) {
            items.push(new SkillTreeItem('本地技能', vscode.TreeItemCollapsibleState.Collapsed));
        }

        return items;
    }
}