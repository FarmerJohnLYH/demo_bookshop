import { Button, InputNumber, message, Modal, Table } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useEffect, useState } from 'react'

interface Book {
  ID: number
  title: string
  author: string
  price: number
  stock: number
}

const CustomerPanel = () => {
  const [books, setBooks] = useState<Book[]>([])
  const [selectedBook, setSelectedBook] = useState<Book | null>(null)
  const [orderQuantity, setOrderQuantity] = useState<number>(1)
  const [isModalVisible, setIsModalVisible] = useState<boolean>(false)

  const fetchBooks = async () => {
    try {
      const response = await fetch('http://localhost:8080/books')
      const data = await response.json()
      setBooks(data)
    } catch (error) {
      message.error('获取图书列表失败')
    }
  }

  useEffect(() => {
    fetchBooks()
  }, [])

  const showOrderModal = (book: Book) => {
    setSelectedBook(book)
    setOrderQuantity(1)
    setIsModalVisible(true)
  }

  const handleOrder = async () => {
    if (!selectedBook) return

    try {
      const response = await fetch('http://localhost:8080/order', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          book_id: selectedBook.ID,
          quantity: orderQuantity,
        }),
      })

      if (response.ok) {
        message.success('下单成功')
        setIsModalVisible(false)
        fetchBooks()
      } else {
        const data = await response.json()
        message.error(data.error || '下单失败')
      }
    } catch (error) {
      message.error('操作失败')
    }
  }

  const columns: ColumnsType<Book> = [
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
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Button
          type="primary"
          onClick={() => showOrderModal(record)}
          disabled={record.stock <= 0}
        >
          购买
        </Button>
      ),
    },
  ]

  return (
    <div style={{ padding: '20px' }}>
      <h2>图书列表</h2>
      <Table
        columns={columns}
        dataSource={books}
        rowKey="ID"
        pagination={false}
      />

      <Modal
        title="确认订单"
        open={isModalVisible}
        onOk={handleOrder}
        onCancel={() => setIsModalVisible(false)}
      >
        {selectedBook && (
          <div>
            <p>书名：{selectedBook.title}</p>
            <p>单价：¥{selectedBook.price.toFixed(2)}</p>
            <p>库存：{selectedBook.stock}</p>
            <div style={{ marginTop: '16px' }}>
              <span style={{ marginRight: '8px' }}>购买数量：</span>
              <InputNumber
                min={1}
                max={selectedBook.stock}
                value={orderQuantity}
                onChange={(value) => setOrderQuantity(value || 1)}
              />
            </div>
            <p style={{ marginTop: '16px' }}>
              总价：¥{((selectedBook.price * orderQuantity) || 0).toFixed(2)}
            </p>
          </div>
        )}
      </Modal>
    </div>
  )
}

export default CustomerPanel