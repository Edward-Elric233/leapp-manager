import React, { useState } from 'react';
import { Button, Form, Header, Segment } from 'semantic-ui-react';
import { API, showError, showSuccess } from '../../helpers';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';

const AddTask = () => {
  const originInputs = {
    task_name: '',
    ip: '',
    port: 22,
  };
  const [inputs, setInputs] = useState(originInputs);
  const { task_name, ip, port } = inputs;
  let navigate = useNavigate();

  const handleInputChange = (e, { name, value }) => {
    setInputs((inputs) => ({ ...inputs, [name]: value }));
  };

  const submit = async () => {
    if (inputs.task_name === '' || inputs.ip === '') return;
    const res = await API.post(`/api/task/`, inputs);
    const { success, message } = res.data;
    if (success) {
      showSuccess('升级任务创建成功！');
      //setInputs(originInputs);
      navigate('/task');
    } else {
      showError(message);
    }

  };

  return (
    <>
      <Segment>
        <Header as="h3">创建新升级任务</Header>
        <Form autoComplete="off">
          <Form.Field>
            <Form.Input
              label="任务名称"
              name="task_name"
              placeholder={'请输入任务名称'}
              onChange={handleInputChange}
              value={task_name}
              autoComplete="off"
              required
            />
          </Form.Field>
          <Form.Field>
            <Form.Input
              label="IP地址"
              name="ip"
              placeholder={'请输入IP地址，例如127.0.0.1'}
              onChange={handleInputChange}
              value={ip}
              autoComplete="off"
            />
          </Form.Field>
          <Form.Field>
            <Form.Input
              label="SSH端口号"
              name="port"
              placeholder={'请输入SSH端口号'}
              onChange={handleInputChange}
              value={port}
              autoComplete="off"
              required
            />
          </Form.Field>
          <Button type={'submit'} onClick={submit}>
            提交
          </Button>
        </Form>
      </Segment>
    </>
  );
};

export default AddTask;
