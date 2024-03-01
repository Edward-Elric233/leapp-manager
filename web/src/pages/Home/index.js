import React, {useContext, useEffect, useState} from 'react';
import { Card, Grid, Header, Segment } from 'semantic-ui-react';
import { API, showError, showNotice, timestamp2string } from '../../helpers';
import { StatusContext } from '../../context/Status';

const Home = () => {
  const [statusState, statusDispatch] = useContext(StatusContext);
  const [tasks, setTasks] = useState(null)
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
        <>
          <Segment>
            <Header as='h3'>Hi, 欢迎使用{statusState?.status?.system_name} !</Header>
          </Segment>
        </>
      )}
    </>
  );
};

export default Home;
