import { Route, Routes } from 'react-router'
import Home from './components/pages/home/Home'
import Article from './components/pages/article/Article'
import Simulator from './components/pages/simulator/Simulator'

function App() {
    return (
        <Routes>
            <Route path='/' element={<Home/>} />
            <Route path='/article/:id' element={<Article/>} />
            <Route path='/simulator' element={<Simulator/>} />
            <Route path='*' element={'Not Found'} />
        </Routes>
    )
}

export default App