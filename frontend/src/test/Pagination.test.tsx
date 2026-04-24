import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { Pagination } from '../components/ui/Pagination';

describe('Pagination', () => {
  it('renders nothing when only one page', () => {
    const { container } = render(<Pagination page={1} total={5} pageSize={10} onChange={() => {}} />);
    expect(container.innerHTML).toBe('');
  });

  it('renders correct page buttons', () => {
    render(<Pagination page={1} total={50} pageSize={10} onChange={() => {}} />);
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('calls onChange with correct page number', () => {
    const onChange = vi.fn();
    render(<Pagination page={2} total={50} pageSize={10} onChange={onChange} />);
    fireEvent.click(screen.getByText('3'));
    expect(onChange).toHaveBeenCalledWith(3);
  });

  it('disables prev button on first page', () => {
    render(<Pagination page={1} total={50} pageSize={10} onChange={() => {}} />);
    const prevButtons = document.querySelectorAll('button');
    expect(prevButtons[0]).toBeDisabled();
  });

  it('disables next button on last page', () => {
    render(<Pagination page={5} total={50} pageSize={10} onChange={() => {}} />);
    const buttons = document.querySelectorAll('button');
    const lastBtn = buttons[buttons.length - 1];
    expect(lastBtn).toBeDisabled();
  });
});
