import { BrowserRouter, Route, Routes } from "react-router-dom";
import { MainLayout } from "./layouts/MainLayout";
// import { AdminLayout } from "./layouts/AdminLayout";
// import { IDELayout } from "./layouts/IDELayout";
import { Home } from "./pages/Home";
import { Search } from "./pages/Search";
import { Detail } from "./pages/Detail";
import { Admin } from "./pages/Admin";
import { IDE } from "./pages/IDE";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<MainLayout />}>
          <Route path="/" element={<Home />} />
          <Route path="/search" element={<Search />} />
          <Route path="/skill/:id" element={<Detail />} />
        </Route>
        
        {/* We'll define specific layouts for these paths since they differ from MainLayout */}
        <Route path="/admin" element={<Admin />} />
        <Route path="/ide" element={<IDE />} />
      </Routes>
    </BrowserRouter>
  );
}
