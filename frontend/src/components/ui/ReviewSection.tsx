import { useState } from "react";
import { StarRating } from "./StarRating";
import { useAuth } from "../../stores/AuthContext";
import { Link } from "react-router-dom";
import type { Review } from "../../types";

interface Props {
  reviews: Review[];
  skillId: number;
}

export function ReviewSection({ reviews, skillId: _skillId }: Props) {
  const { isAuthenticated } = useAuth();
  const [rating, setRating] = useState(0);
  const [comment, setComment] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [localReviews, setLocalReviews] = useState(reviews);

  const handleSubmit = async () => {
    if (!rating || !comment.trim()) return;
    setSubmitting(true);
    try {
      await new Promise(r => setTimeout(r, 500));
      const newReview: Review = {
        id: Date.now(),
        userId: 0,
        userName: '您',
        userAvatar: '',
        rating,
        comment: comment.trim(),
        createdAt: '刚刚',
      };
      setLocalReviews(prev => [newReview, ...prev]);
      setRating(0);
      setComment("");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col gap-6">
      <h3 className="text-lg font-bold text-slate-900 flex items-center gap-2">
        <span className="material-symbols-outlined text-brand-600">rate_review</span>
        评分与评论 ({localReviews.length})
      </h3>

      {isAuthenticated ? (
        <div className="card p-5">
          <h4 className="font-medium text-slate-900 mb-3">撰写评论</h4>
          <div className="flex flex-col gap-3">
            <div className="flex items-center gap-2">
              <span className="text-sm text-slate-600">评分:</span>
              <StarRating rating={rating} size="md" interactive onChange={setRating} />
            </div>
            <textarea
              value={comment}
              onChange={e => setComment(e.target.value)}
              className="input-field h-24 resize-none"
              placeholder="分享您的使用体验..."
            />
            <div className="flex justify-end">
              <button
                onClick={handleSubmit}
                disabled={!rating || !comment.trim() || submitting}
                className="btn-primary text-sm"
              >
                {submitting ? '提交中...' : '提交评论'}
              </button>
            </div>
          </div>
        </div>
      ) : (
        <div className="card p-5 text-center">
          <p className="text-sm text-slate-500">
            <Link to="/login" className="text-brand-600 font-medium hover:text-brand-700">登录</Link> 后即可发表评论
          </p>
        </div>
      )}

      <div className="flex flex-col gap-4">
        {localReviews.length === 0 ? (
          <p className="text-sm text-slate-500 text-center py-8">暂无评论，成为第一个评价的人吧</p>
        ) : (
          localReviews.map(review => (
            <div key={review.id} className="card p-4">
              <div className="flex items-start justify-between mb-2">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-slate-200 overflow-hidden">
                    {review.userAvatar ? (
                      <img src={review.userAvatar} alt="" className="w-full h-full object-cover" />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-xs font-bold text-slate-500 bg-slate-300">
                        {review.userName.charAt(0)}
                      </div>
                    )}
                  </div>
                  <div>
                    <p className="text-sm font-medium text-slate-900">{review.userName}</p>
                    <span className="text-xs text-slate-400">{review.createdAt}</span>
                  </div>
                </div>
                <StarRating rating={review.rating} size="sm" />
              </div>
              <p className="text-sm text-slate-700 leading-relaxed">{review.comment}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
