import React, { useEffect, useState } from 'react';
import {
  Button,
  Form,
  Label,
  Pagination,
  Popup,
  Table,
} from 'semantic-ui-react';
import { Link } from 'react-router-dom';
import { API, showError, showSuccess } from '../helpers';

import { ITEMS_PER_PAGE } from '../constants';

const TasksTable = () => {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activePage, setActivePage] = useState(1);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [searching, setSearching] = useState(false);

  const loadTasks = async (startIdx) => {
    const data = [{task_name: "测试机器1", ip: "192.168.1.1", status: 1}]
    setTasks(prevState => {
      return [...prevState, ...data];
    });
    /*
    const res = await API.get(`/api/task/?p=${startIdx}`);
    const { success, message, data } = res.data;
    if (success) {
      if (startIdx === 0) {
        setTasks(data);
      } else {
        setTasks(prevState => {
          return [...prevState, ...data];
        });
      }
    } else {
      showError(message);
    }
    */
    setLoading(false);
  };

  const onPaginationChange = (e, { activePage }) => {
    (async () => {
      if (activePage === Math.ceil(tasks.length / ITEMS_PER_PAGE) + 1) {
        // In this case we have to load more data and then append them.
        await loadTasks(activePage - 1);
      }
      setActivePage(activePage);
    })();
  };

  useEffect(() => {
    loadTasks(0)
        .then()
        .catch((reason) => {
          showError(reason);
        });
  }, []);

  const manageTask = (username, action, idx) => {
    (async () => {
      const res = await API.post('/api/task/manage', {
        username,
        action,
      });
      const { success, message } = res.data;
      if (success) {
        showSuccess('操作完成！');
        let user = res.data.data;
        let newUsers = [...tasks];
        let realIdx = (activePage - 1) * ITEMS_PER_PAGE + idx;
        if (action === 'delete') {
          newUsers[realIdx].deleted = true;
        } else {
          newUsers[realIdx].status = user.status;
          newUsers[realIdx].role = user.role;
        }
        setTasks(newUsers);
      } else {
        showError(message);
      }
    })();
  };

  const renderStatus = (status) => {
    switch (status) {
      case 1:
        return <Label basic>正在升级</Label>;
      case 2:
        return (
            <Label basic color='green'>
              升级成功
            </Label>
        );
      case 3:
        return (
            <Label basic color='red'>
              升级失败
            </Label>
        );
      default:
        return (
            <Label basic color='grey'>
              未开始
            </Label>
        );
    }
  };

  const searchTasks = async () => {
    if (searchKeyword === '') {
      // if keyword is blank, load files instead.
      await loadTasks(0);
      setActivePage(1);
      return;
    }
    setSearching(true);
    const res = await API.get(`/api/task/search?keyword=${searchKeyword}`);
    const { success, message, data } = res.data;
    if (success) {
      setTasks(data);
      setActivePage(1);
    } else {
      showError(message);
    }
    setSearching(false);
  };

  const handleKeywordChange = async (e, { value }) => {
    setSearchKeyword(value.trim());
  };

  const sortTask = (key) => {
    if (tasks.length === 0) return;
    setLoading(true);
    setTasks(prevState => {
      let sortedUsers = [...prevState];
      sortedUsers.sort((a, b) => {
        return ('' + a[key]).localeCompare(b[key]);
      });
      if (sortedUsers[0].id === tasks[0].id) {
        sortedUsers.reverse();
      }
      return sortedUsers;
    });
    setLoading(false);
  };

  return (
      <>
        <Form onSubmit={searchTasks}>
          <Form.Input
              icon='search'
              fluid
              iconPosition='left'
              placeholder='需要搜索的任务名称/IP地址...'
              value={searchKeyword}
              loading={searching}
              onChange={handleKeywordChange}
          />
        </Form>

        <Table basic>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell
                  style={{ cursor: 'pointer' }}
                  onClick={() => {
                    sortTask('task_name');
                  }}
              >
                任务名称
              </Table.HeaderCell>
              <Table.HeaderCell
                  style={{ cursor: 'pointer' }}
                  onClick={() => {
                    sortTask('ip');
                  }}
              >
                IP地址
              </Table.HeaderCell>
              <Table.HeaderCell
                  style={{ cursor: 'pointer' }}
                  onClick={() => {
                    sortTask('status');
                  }}
              >
                任务状态
              </Table.HeaderCell>
              <Table.HeaderCell>操作</Table.HeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            {tasks
                .slice(
                    (activePage - 1) * ITEMS_PER_PAGE,
                    activePage * ITEMS_PER_PAGE
                )
                .map((task, idx) => {
                  if (task.deleted) return <></>;
                  return (
                      <Table.Row key={task.id}>
                        <Table.Cell>{task.task_name}</Table.Cell>
                        <Table.Cell>{task.ip}</Table.Cell>
                        <Table.Cell>{renderStatus(task.status)}</Table.Cell>
                        <Table.Cell>
                          <div>
                            {
                              //TODO: 操作： 对于不是正在升级的任务， 可以从数据库中删除
                            }
                            <Popup
                                trigger={
                                  <Button size='small' negative>
                                    删除
                                  </Button>
                                }
                                on='click'
                                flowing
                                hoverable
                            >
                            </Popup>
                            <Button
                                size={'small'}
                                onClick={() => {
                                  manageTask(
                                      task.username,
                                      task.status === 1 ? 'disable' : 'enable',
                                      idx
                                  );
                                }}
                            >
                              {task.status === 1 ? '停止' : '开始'}
                            </Button>
                          </div>
                        </Table.Cell>
                      </Table.Row>
                  );
                })}
          </Table.Body>

          <Table.Footer>
            <Table.Row>
              <Table.HeaderCell colSpan='6'>
                <Button size='small' as={Link} to='/task/add' loading={loading}>
                  新建任务
                </Button>
                <Pagination
                    floated='right'
                    activePage={activePage}
                    onPageChange={onPaginationChange}
                    size='small'
                    siblingRange={1}
                    totalPages={
                        Math.ceil(tasks.length / ITEMS_PER_PAGE) +
                        (tasks.length % ITEMS_PER_PAGE === 0 ? 1 : 0)
                    }
                />
              </Table.HeaderCell>
            </Table.Row>
          </Table.Footer>
        </Table>
      </>
  );
};

export default TasksTable;
