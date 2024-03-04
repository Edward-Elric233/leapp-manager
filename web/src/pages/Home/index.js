import React, {useContext, useEffect, useState} from 'react';
import { Card, Grid, Header, Segment, Container, Icon, Button } from 'semantic-ui-react';
import { API, showError, showNotice, timestamp2string } from '../../helpers';
import { StatusContext } from '../../context/Status';
import {UserContext} from "../../context/User";
import TasksTable from "../../components/TasksTable";

const Home = () => {
  const [statusState, statusDispatch] = useContext(StatusContext);
  const [userState, userDispatch] = useContext(UserContext);
  const homePageLink = localStorage.getItem('home_page_link') || '';

  const displayNotice = async () => {
    const res = await API.get('/api/notice');
    const { success, message, data } = res.data;
    if (success) {
      let oldNotice = localStorage.getItem('notice');
      if (data !== oldNotice && data !== '') {
        showNotice(data);
        localStorage.setItem('notice', data);
      }
    } else {
      showError(message);
    }
  };

  const getStartTimeString = () => {
    const timestamp = statusState?.status?.start_time;
    return timestamp2string(timestamp);
  };

  useEffect(() => {
    displayNotice().then();
  }, []);
  return (
      <>
        {homePageLink !== '' ? (
            <>
              <iframe
                  src={homePageLink}
                  style={{ width: '100%', height: '100vh', border: 'none' }}
              />
            </>
        ) : (
            <Container>
              <Segment padded='very' textAlign='center' style={{ minHeight: '80vh' }}>
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
        )}
      </>
  );
};

export default Home;
