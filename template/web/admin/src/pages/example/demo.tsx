import { ProCard } from '@ant-design/pro-components';
import { Button, Space, Table, Tag } from 'antd';
import { PlusOutlined } from '@ant-design/icons';

const columns = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (status: string) => (
      <Tag color={status === 'running' ? 'green' : 'default'}>{status}</Tag>
    ),
  },
  {
    title: '更新时间',
    dataIndex: 'updatedAt',
    key: 'updatedAt',
  },
];

const mockData = [
  { key: '1', name: '示例条目 A', status: 'running', updatedAt: '2026-01-01 12:00:00' },
  { key: '2', name: '示例条目 B', status: 'stopped', updatedAt: '2026-01-02 08:30:00' },
];

function ExampleDemo() {
  return (
    <ProCard
      title="示例页面"
      extra={
        <Button type="primary" icon={<PlusOutlined />}>
          新增
        </Button>
      }
    >
      <Space direction="vertical" style={{ width: '100%' }} size="middle">
        <Table columns={columns} dataSource={mockData} pagination={false} />
      </Space>
    </ProCard>
  );
}

export default ExampleDemo;
