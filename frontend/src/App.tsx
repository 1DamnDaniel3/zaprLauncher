import './App.css';
import icon from '../../build/windows/icon.ico'
import { Loader } from './shared/Loader/Loader';

function App() {


    return (
        <div id="App">
            <img className='icon' src={icon}></img>
            <Loader />
        </div>
    )
}

export default App
