import React from 'react';
import {ThemeProvider} from "@components/organisms";

import "antd/dist/reset.css";


export const wrapRootElement = ({ element }: { element: React.ReactNode }) => {
    return (
        <ThemeProvider>
            {element}
        </ThemeProvider>
    )
}
