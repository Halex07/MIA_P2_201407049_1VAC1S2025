import { useState } from "react";
import { useParams } from 'react-router-dom';
import { useNavigate } from "react-router-dom"
import partIMG from '../iconos/part.png';
import "../StyleSheets/fondo.css"

export default function Partitions({newIp="0.0.0.0"}){
    const { id } = useParams()
    const [ particiones, setParticiones ] = useState([]);
    const navigate = useNavigate()
    
    useState(()=>{
        fetch(`http://${newIp}:8080/particiones`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify(id)
        })
        .then(Response => Response.json())
        .then(rawData => {console.log(rawData); setParticiones(rawData);})
        .catch(error => {
            console.error('Error en la solicitud Fetch:', error);
            
        });
    }, [])

    const onClick = (particion) => {
        console.log("click",particion)
        navigate(`/login/${id}/${particion}`)
    }

    return(
        <div className='body'>
            <div>&nbsp;&nbsp;&nbsp;</div>
            <div style={{display:"flex", flexDirection:"row", justifyContent: "center"}}><h1>Particion en Disco {id} </h1></div>
            <div style={{display:"flex", flexDirection:"row", justifyContent: "center"}}>
                {particiones && particiones.length > 0 ? (
                    particiones.map((particion, index) => {
                        return (
                            <div key={index} style={{
                                display: "flex",
                                flexDirection: "column", 
                                alignItems: "center", 
                                maxWidth: "100px",
                                margin: "10px"
                                }}
                                onClick={() => onClick(particion)}
                            >
                                <img src={partIMG} alt="part" style={{width: "100px"}} />
                                {particion}
                            </div>
                        )
                    })
                ):(
                    <div>No existen particiones disponibles</div>
                )}
            </div> 
        </div> 
    );
}