import { Component, type ErrorInfo, type ReactNode } from "react";

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("ErrorBoundary caught an error:", error, errorInfo);
  }

  handleRetry = () => {
    this.setState({ hasError: false, error: null });
  };

  render() {
    if (this.state.hasError) {
      return (
        <div
          role="alert"
          className="min-h-screen flex items-center justify-center bg-slate-50 p-8"
        >
          <div className="bg-white rounded-xl border border-red-200 shadow-card p-8 max-w-md w-full text-center flex flex-col items-center gap-4">
            <span className="material-symbols-outlined text-[64px] text-red-400">error</span>
            <h1 className="text-xl font-bold text-slate-900">Something went wrong</h1>
            <p className="text-sm text-slate-500">
              An unexpected error occurred. Please try again or contact support if the problem persists.
            </p>
            {this.state.error && (
              <details className="w-full text-left">
                <summary className="text-xs text-slate-400 cursor-pointer hover:text-slate-600">
                  Error details
                </summary>
                <pre className="mt-2 text-xs text-red-600 bg-red-50 p-3 rounded border border-red-100 overflow-x-auto">
                  {this.state.error.message}
                </pre>
              </details>
            )}
            <button
              onClick={this.handleRetry}
              className="btn-primary mt-2"
            >
              <span className="material-symbols-outlined text-[18px]">refresh</span>
              Try Again
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
