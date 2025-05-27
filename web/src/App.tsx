import { createTheme, ThemeProvider } from '@mui/material';
import './App.css';
import Home from './Home';



const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

function App() {

  return (
    <ThemeProvider theme={darkTheme}>
      <Home />
    </ThemeProvider>
  )
}

export default App
