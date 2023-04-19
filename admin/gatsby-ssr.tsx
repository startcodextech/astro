import * as React from "react"

export const onRenderBody = ({ setHeadComponents }) => {
    setHeadComponents([
        <link key="googleapis" rel="preconnect" href="https://fonts.googleapis.com"/>,
        <link key="gstatic" rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous"/>,
        <link key="redhat" href="https://fonts.googleapis.com/css2?family=Red+Hat+Text:wght@300;400;500;700&display=swap" rel="stylesheet"/>
    ])
}
