import React, {PropsWithChildren, FC} from "react";
import {ConfigProvider, ThemeConfig} from "antd";
import esEs from 'antd/es/locale/es_ES';

export const defaultTheme: Partial<ThemeConfig> = {
    token: {
        fontFamily: `'Red Hat Text', sans-serif`,
        fontWeightStrong: 700,
    }
}

type Props = PropsWithChildren<{}>

const ThemeProvider: FC<Props> = (props) => {
    const {children} = props;
    return (
        <ConfigProvider theme={defaultTheme} locale={esEs}>
            {children}
        </ConfigProvider>
    )
}

export default ThemeProvider;
