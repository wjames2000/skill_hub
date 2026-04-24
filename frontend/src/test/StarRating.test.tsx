import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { StarRating } from '../components/ui/StarRating';

describe('StarRating', () => {
  it('renders the correct number of stars', () => {
    render(<StarRating rating={3.5} />);
    const stars = document.querySelectorAll('.material-symbols-outlined');
    expect(stars.length).toBe(5);
  });

  it('displays the rating value', () => {
    render(<StarRating rating={4.2} />);
    expect(screen.getByText('4.2')).toBeInTheDocument();
  });

  it('calls onChange when star is clicked in interactive mode', () => {
    const onChange = vi.fn();
    render(<StarRating rating={0} interactive onChange={onChange} />);
    const stars = document.querySelectorAll('button');
    fireEvent.click(stars[2]);
    expect(onChange).toHaveBeenCalledWith(3);
  });

  it('disables click when not interactive', () => {
    const onChange = vi.fn();
    render(<StarRating rating={4} onChange={onChange} />);
    const stars = document.querySelectorAll('button');
    fireEvent.click(stars[0]);
    expect(onChange).not.toHaveBeenCalled();
  });

  it('renders correct size class', () => {
    const { container } = render(<StarRating rating={3} size="lg" />);
    expect(container.querySelector('button')?.className).toContain('text-[24px]');
  });
});
