import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { SkillCard } from '../components/ui/SkillCard';
import { LanguageProvider } from '../stores/LanguageContext';
import type { Skill } from '../types';

const mockSkill: Skill = {
  id: 1,
  title: 'Test Skill',
  description: 'A test skill description',
  zhDescription: '',
  enDescription: '',
  author: 'Test Author',
  icon: 'terminal',
  iconColor: 'text-green-600',
  iconBg: 'bg-green-50',
  tags: ['test', 'demo'],
  category: 'Testing',
  version: 'v1.0.0',
  rating: 4.5,
  downloads: 5000,
  installCount: 5000,
  source: 'official',
  safe: true,
  createdAt: '2024-01-01',
  updatedAt: '2024-06-01',
};

function renderWithRouter(ui: React.ReactElement) {
  return render(<BrowserRouter><LanguageProvider>{ui}</LanguageProvider></BrowserRouter>);
}

describe('SkillCard', () => {
  it('renders skill title', () => {
    renderWithRouter(<SkillCard skill={mockSkill} />);
    expect(screen.getByText('Test Skill')).toBeInTheDocument();
  });

  it('renders author name', () => {
    renderWithRouter(<SkillCard skill={mockSkill} />);
    expect(screen.getByText(/Test Author/)).toBeInTheDocument();
  });

  it('renders rating', () => {
    renderWithRouter(<SkillCard skill={mockSkill} />);
    expect(screen.getByText('4.5')).toBeInTheDocument();
  });

  it('renders tags', () => {
    renderWithRouter(<SkillCard skill={mockSkill} />);
    expect(screen.getByText('test')).toBeInTheDocument();
    expect(screen.getByText('demo')).toBeInTheDocument();
  });

  it('links to skill detail page', () => {
    renderWithRouter(<SkillCard skill={mockSkill} />);
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/skill/1');
  });

  it('shows safe badge when skill is safe', () => {
    renderWithRouter(<SkillCard skill={mockSkill} />);
    expect(screen.getByText('verified')).toBeInTheDocument();
  });

  it('shows match score when prop is set', () => {
    const skillWithMatch = { ...mockSkill, matchScore: 92 };
    renderWithRouter(<SkillCard skill={skillWithMatch} showMatchScore />);
    expect(screen.getByText('92%')).toBeInTheDocument();
  });
});
