import { Routes, Route, HashRouter, Link } from 'react-router-dom'
import { useState } from 'react';

import hacker from '../iconos/hacker.gif';
import Comandos from '../Paginas/comandos';
import Discos from '../Paginas/discos';
import Partitions from '../Paginas/partition';
import Login from '../Paginas/login';
import Explorer from '../Paginas/explorador';
import Reportes from '../Paginas/reportes';

export default function Navegador(){
    const [ ip, setIP ] = useState("0.0.0.0")
    
    const handleChange = (e) => {
        console.log(e.target.value)
        setIP(e.target.value)
    }
    
    const logOut = (e) => {
        e.preventDefault()
        
        fetch(`http://${ip}:8080/logout`)
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData);  
            if (rawData === 0){
                alert('sesion cerrada')
                window.location.href = '#/Discos';
            }else{
                alert('No hay sesion abierta')
            }
        }) 
        .catch(error => {
            console.error('Error en la solicitud Fetch:', error);
          
        });
    };

    const limpiar = (e) => {
        e.preventDefault()
        console.log("limpiando")
        fetch(`http://${ip}:8080/limpiar`)
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData); 
            if (rawData === 1){
                alert('Discos y reportes eliminados')
                window.location.href = '#/Comandos';
            }else{
                alert('Error al eliminar archiovs')
            }
        }) 
    }

    return(
        <HashRouter>
            <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
             
                <div id="espacio">&nbsp;&nbsp;&nbsp;</div>
                
                <div className="conteiner-fluid"> 
                    <img src={hacker} alt="" width="64" height="64" className="d-inline-block align-text-top"></img>
                </div>

                <div className="conteiner"> 
                   
                    <div className="container-fluid">
                        <a className="navbar-brand" type="submit" >
                            ARCHIVOS PROYECTO 2            
                        </a>
                      
                        <div className="collapse navbar-collapse" id="navbarColor02">
              
                            <ul className="navbar-nav me-auto">
                              
                                <li className="nav-item">
                               
                                    <Link className="nav-link active" to="/Comandos">Comandos</Link>
                                </li>

                                <li className="nav-item">
                             
                                    <Link className="nav-link" to="/Discos">Explorador</Link>
                                </li>

                                <li className="nav-item">
                                    <button onClick={logOut} className="nav-link">Logout</button>
                                </li>

                                <li className="nav-item">
                                    <Link className="nav-link" to="/Reportes">Reportes</Link>
                                </li>

                                <li className="nav-item">
                                    <button onClick={limpiar} className="nav-link">Limpiar</button>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
                <input className="form-control me-2 mx-auto" style={{ maxWidth: "200px" }} placeholder="IP" onChange={handleChange}/>
                <div id="espacio">&nbsp;</div>
            </nav> 
            
            <Routes>
                <Route path="/" element ={<Comandos newIp={ip}/>}/>
                <Route path="/Comandos" element ={<Comandos newIp={ip}/>}/> 
                <Route path="/Discos" element ={<Discos newIp={ip}/>}/> 
                <Route path="/Disco/:id" element ={<Partitions newIp={ip}/>}/> 
                <Route path="/Login/:disk/:part" element ={<Login newIp={ip}/>}/>
                <Route path="/Explorador/:id" element ={<Explorer newIp={ip}/>}/>
                <Route path="/Reportes" element ={<Reportes newIp={ip}/>}/>                 
            </Routes>
        </HashRouter>
    );
}
