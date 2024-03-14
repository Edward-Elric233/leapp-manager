import React, { useState, useEffect } from 'react';
import { Button, Form, TextArea, Segment, Container, Header} from 'semantic-ui-react';
import {API, showError, showSuccess, showWarning} from '../../helpers';
import {useParams} from "react-router-dom";
import Loading from "../../components/Loading";
import NotFound from "../NotFound";
import {taskConstants} from "../../constants";

const TaskDetails = () => {
    const [task, setTask] = useState(null);
    const [loading, setLoading] = useState(true);
    const params = useParams()
    const id = params.id

    const fetchTask = async () => {
        try {
            const res = await API.get(`/api/task/${id}`);
            const { success, data, message } = res.data;
            if (success) {
                setTask(data);
            } else {
                showError(message);
            }
            setLoading(false);
        } catch (error) {
            showError(error.message);
            setLoading(false);
        }
    };
    useEffect(() => {
        fetchTask()
            .then()
            .catch(reason => showError(reason))
    }, [id]);

    const handleStartUpgrade = async () => {
        try {
            const res = await API.post(`/api/task/${task.id}`);
            const { success, message, data } = res.data;
            if (success) {
                setTask(data);
            } else {
                showError(message);
            }
        } catch (error) {
            showError(error.message);
        }
    };

    if (loading) {
        return <Loading/>;
    }

    if (!task) {
        return <NotFound/>;
    }

    const {taskStatus} = taskConstants
    let buttonName = ""
    //console.log("任务状态:", task.status, taskStatus[task.status])
    switch (taskStatus[task.status]) {
        case "init":
            buttonName = "开始升级";
            break;
        case "run":
            buttonName = "升级中";
            break;
        default:
            buttonName = "已结束"
    }

    return (
        <Container>
            <Header as='h1' attached='top' block>
                升级任务详情
            </Header>
            <Segment attached>
                <p><strong>任务名称：</strong>{task.task_name}</p>
                <p><strong>IP信息：</strong>{task.ip}</p>
                <Button
                    onClick={() => handleStartUpgrade()}
                    disabled={taskStatus[task.status] !== 'init'}
                    color='green'
                >
                    {buttonName}
                </Button>
            </Segment>
            <Form>
                <TextArea
                    placeholder='升级详情将显示在这里...'
                    value={task.info}
                    style={{ minHeight: '50vh' }}
                    readOnly
                />
            </Form>
        </Container>
    );
};

export default TaskDetails;