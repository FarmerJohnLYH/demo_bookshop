import { SettingOutlined, ShopOutlined } from '@ant-design/icons'
import { Tabs } from 'antd'
import './App.css'
import AdminPanel from './components/AdminPanel'
import CustomerPanel from './components/CustomerPanel'

function App() {
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

export default App
