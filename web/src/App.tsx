import { AppBar, createTheme, ThemeProvider, Toolbar, Typography } from '@mui/material';
import axios from 'axios';
import { BrowserRouter } from 'react-router';
import './App.css';
import { AppRoutes } from './Routes';



const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

function App() {
  // axios.defaults.baseURL = 'https://mn6v9itpxb.us-east-2.awsapprunner.com/';
  axios.defaults.baseURL = 'http://localhost:8080/v0/';

  return (
    <ThemeProvider theme={darkTheme}>
      <BrowserRouter>
        <AppBar>
          <Toolbar>
            <Typography
              variant="h6"
              noWrap
            >
              Rikon Parts Stock Alert
            </Typography>
          </Toolbar>
        </AppBar>
        <AppRoutes />
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
