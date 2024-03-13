import React, {useContext, useEffect, useState} from 'react';
import { API, showError, showNotice, timestamp2string } from '../../helpers';
import { StatusContext } from '../../context/Status';
import {UserContext} from "../../context/User";
import TasksTable from "../../components/TasksTable";
import Task from "../Task";

const Home = () => {
  const [statusState, statusDispatch] = useContext(StatusContext);
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
            <Task/>
        )}
      </>
  );
};

export default Home;
