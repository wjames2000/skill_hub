import { BrowserRouter, Route, Routes } from "react-router-dom";
import { lazy, Suspense } from "react";
import { MainLayout } from "./layouts/MainLayout";
import { AuthProvider } from "./stores/AuthContext";
import { LanguageProvider } from "./stores/LanguageContext";
import { LoadingScreen } from "./components/ui/LoadingScreen";
import { ProtectedRoute } from "./components/ui/ProtectedRoute";
import { ErrorBoundary } from "./components/ui/ErrorBoundary";

const Home = lazy(() => import("./pages/Home").then(m => ({ default: m.Home })));
const Search = lazy(() => import("./pages/Search").then(m => ({ default: m.Search })));
const Detail = lazy(() => import("./pages/Detail").then(m => ({ default: m.Detail })));
const Login = lazy(() => import("./pages/Login").then(m => ({ default: m.Login })));
const Register = lazy(() => import("./pages/Register").then(m => ({ default: m.Register })));
const Profile = lazy(() => import("./pages/Profile").then(m => ({ default: m.Profile })));
const Admin = lazy(() => import("./pages/Admin").then(m => ({ default: m.Admin })));
const IDE = lazy(() => import("./pages/IDE").then(m => ({ default: m.IDE })));

export default function App() {
  return (
    <BrowserRouter>
      <LanguageProvider>
      <AuthProvider>
        <Suspense fallback={<LoadingScreen />}>
          <ErrorBoundary>
          <Routes>
            <Route element={<MainLayout />}>
              <Route path="/" element={<Home />} />
              <Route path="/search" element={<Search />} />
              <Route path="/skill/:id" element={<Detail />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/profile" element={<ProtectedRoute><Profile /></ProtectedRoute>} />
            </Route>
            <Route path="/admin" element={<ProtectedRoute requireAdmin><Admin /></ProtectedRoute>} />
            <Route path="/ide" element={<IDE />} />
          </Routes>
          </ErrorBoundary>
        </Suspense>
      </AuthProvider>
      </LanguageProvider>
    </BrowserRouter>
  );
}
