import React, {useContext, useEffect, useState} from 'react';
import { Card, Grid, Header, Segment, Container, Icon, Button } from 'semantic-ui-react';
import TasksTable from '../../components/TasksTable';
import {UserContext} from "../../context/User";
import {StatusContext} from "../../context/Status";

const Task = () => {
    const [userState, userDispatch] = useContext(UserContext);
    const [statusState, statusDispatch] = useContext(StatusContext);
    //TODO: refactor
    return (<>
        <Container>
            <Segment padded='very' textAlign='center' style={{ minHeight: '90vh' }}>
                <Header as='h2' icon>
                    <img
                        src='/logo.png'
                        alt='logo'
                        style={{ marginRight: '0.75em' }}
                    />
                    Hi, 欢迎使用{statusState?.status?.system_name} !
                    <Header.Subheader>
                        管理和监控您的系统升级任务
                    </Header.Subheader>
                </Header>
                {userState.user ? (
                    <Segment padded>
                        <TasksTable />
                    </Segment>
                ) : (
                    <Segment>
                        <Header icon>
                            <Icon name='settings' />
                            请登录进行系统升级
                        </Header>
                    </Segment>
                )}
            </Segment>
        </Container>
    </>)
};

export default Task;
