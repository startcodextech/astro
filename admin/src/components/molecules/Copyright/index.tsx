import React from 'react'
import {version} from '@package';

import './styled.sass'

const Copyright = () => {
    return (
        <div className="copyright">
            <p><span className="name">ASTRO V{version}</span> | Desarrollado por <a target="_blank" href="https://startcodex.com/astro">Start Codex</a></p>
        </div>
    )
}

export default Copyright
