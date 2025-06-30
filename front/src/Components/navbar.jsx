import React, { useState } from 'react';
import hacker from '../iconos/hacker.gif';
import Comandos from '../Paginas/comandos';
import Explorer from '../Paginas/discos';

export default function Navbar(){
    const [componenteActivo, setComponenteActivo] = useState(<Comandos/>);

    function comandos(idPagina){
        let componente
        if (idPagina === 1){
            componente =  <Comandos/>
        }else if (idPagina === 2){
            componente = <Explorer/>
        }
        setComponenteActivo(componente)
    }

    return(
        <>
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
                                    <a className="nav-link active" type="button" onClick={() => comandos(1)}>Comandos</a>
                                </li>

                                <li className="nav-item">
                                    <a className="nav-link" type="button" onClick={() => comandos(2)}>Explorador</a>
                                </li>

                                <li className="nav-item">
                                    <a className="nav-link" type="submit">Logout</a>
                                </li>

                            </ul>
                        </div>
                    </div>
                </div>
            </nav> 
            {componenteActivo}
        </>
    );
}

