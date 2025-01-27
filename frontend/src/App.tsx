import { SettingOutlined, ShopOutlined } from '@ant-design/icons'
import { Tabs } from 'antd'
import { Navigate, Route, BrowserRouter as Router, Routes } from 'react-router-dom'
import './App.css'
import AdminPanel from './components/AdminPanel'
import CustomerPanel from './components/CustomerPanel'
import LoginPage from './components/LoginPage'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem('token')
  return token ? <>{children}</> : <Navigate to="/login" />
}

function MainContent() {
  return (
    <div className="app-container">
      <h1>简易书店系统</h1>
      <Tabs
        defaultActiveKey="customer"
        items={[
          {
            key: 'customer',
            label: (
              <span>
                <ShopOutlined />
                客户端
              </span>
            ),
            children: <CustomerPanel />,
          },
          {
            key: 'admin',
            label: (
              <span>
                <SettingOutlined />
                管理端
              </span>
            ),
            children: <AdminPanel />,
          },
        ]}
      />
    </div>
  )
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <MainContent />
            </PrivateRoute>
          }
        />
      </Routes>
    </Router>
  )
}

export default App
