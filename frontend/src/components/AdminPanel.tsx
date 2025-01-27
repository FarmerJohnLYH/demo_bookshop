import { Button, Form, Input, InputNumber, Table, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useEffect, useState } from 'react'

interface Book {
  ID: number
  title: string
  author: string
  price: number
  stock: number
}

const AdminPanel = () => {
  const [books, setBooks] = useState<Book[]>([])
  const [form] = Form.useForm()

  /**
   * 获取所有图书信息
   * @description 从后端API获取图书列表数据
   * @returns {Promise<void>}
   * @throws {Error} 当API请求失败时抛出错误
   */
  const fetchBooks = async () => {
    try {
      const response = await fetch('http://localhost:8080/books')  // 使用的是 get 方法
      const data = await response.json()
      setBooks(data)
    } catch (error) {
      message.error('获取图书列表失败')
    }
  }

  useEffect(() => {
    fetchBooks()
  }, [])

  const columns: ColumnsType<Book> = [
    {
      title: 'ID',
      dataIndex: 'ID',
      key: 'id',
    },
    {
      title: '书名',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: '作者',
      dataIndex: 'author',
      key: 'author',
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => `¥${price.toFixed(2)}`,
    },
    {
      title: '库存',
      dataIndex: 'stock',
      key: 'stock',
    },
  ]

  const onFinish = async (values: any) => {
    try {
      const response = await fetch('http://localhost:8080/admin/stock/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(values),
      })

      if (response.ok) {
        message.success('库存更新成功')
        form.resetFields()
        fetchBooks()
      } else {
        message.error('库存更新失败')
      }
    } catch (error) {
      message.error('操作失败')
    }
  }

  return (
    <div style={{ padding: '20px' }}>
      <h2>图书库存管理</h2>
      <div style={{ marginBottom: '20px' }}>
        <Button
          type="primary"
          danger
          onClick={async () => {
            try {
              const response = await fetch('http://localhost:8080/admin/reset-books', {
                method: 'POST',
              })
              if (response.ok) {
                message.success('图书数据已重置')
                fetchBooks()
              } else {
                message.error('重置数据失败')
              }
            } catch (error) {
              message.error('操作失败')
            }
          }}
          style={{ marginRight: '16px' }}
        >
          重置图书数据
        </Button>
      </div>
      <Form
        form={form}
        name="bookForm"
        onFinish={onFinish}
        layout="inline"
        style={{ marginBottom: '20px' }}
      >
        <Form.Item
          name="title"
          rules={[{ required: true, message: '请输入书名' }]}
        >
          <Input placeholder="书名" />
        </Form.Item>

        <Form.Item
          name="author"
          rules={[{ required: true, message: '请输入作者' }]}
        >
          <Input placeholder="作者" />
        </Form.Item>

        <Form.Item
          name="price"
          rules={[{ required: true, message: '请输入价格' }]}
        >
          <InputNumber
            placeholder="价格"
            min={0}
            step={0.01}
            style={{ width: '100px' }}
          />
        </Form.Item>

        <Form.Item
          name="stock"
          rules={[{ required: true, message: '请输入库存数量' }]}
        >
          <InputNumber
            placeholder="库存"
            min={1}
            style={{ width: '100px' }}
          />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit">
            添加/更新库存
          </Button>
        </Form.Item>
      </Form>

      <Table
        columns={columns}
        dataSource={books}
        rowKey="ID"
        pagination={false}
      />
    </div>
  )
}

export default AdminPanel