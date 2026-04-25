interface Props {
  rating: number;
  size?: 'sm' | 'md' | 'lg';
  interactive?: boolean;
  onChange?: (rating: number) => void;
}

export function StarRating({ rating, size = 'md', interactive, onChange }: Props) {
  const sizeClasses = { sm: 'text-[14px]', md: 'text-[18px]', lg: 'text-[24px]' };
  const fullStars = Math.floor(rating);
  const hasHalf = rating - fullStars >= 0.5;

  return (
    <div className="flex items-center gap-0.5" role={interactive ? 'radiogroup' : undefined} aria-label={interactive ? 'Rating' : undefined}>
      {[1, 2, 3, 4, 5].map((star) => {
        const filled = star <= fullStars;
        const half = star === fullStars + 1 && hasHalf;
        return (
          <button
            key={star}
            type="button"
            disabled={!interactive}
            aria-label={interactive ? `Rate ${star} star${star > 1 ? 's' : ''}` : undefined}
            onClick={() => interactive && onChange?.(star)}
            className={`${interactive ? 'cursor-pointer hover:scale-110' : 'cursor-default'} transition-transform ${sizeClasses[size]}`}
          >
            <span className={`material-symbols-outlined ${filled || half ? 'fill' : ''} text-amber-500`}>
              {half ? 'star_half' : 'star'}
            </span>
          </button>
        );
      })}
      <span className={`ml-1 font-medium text-slate-600 ${size === 'sm' ? 'text-xs' : 'text-sm'}`}>
        {rating.toFixed(1)}
      </span>
    </div>
  );
}
