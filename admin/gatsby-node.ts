import {GatsbyNode} from 'gatsby';
import path from 'path';

export const onCreateWebpackConfig: GatsbyNode['onCreateWebpackConfig'] = (props) => {
    const {actions} = props;

    actions.setWebpackConfig({
        resolve: {
            alias: {
                '@components': path.resolve(__dirname, 'src/components'),
                '@assets': path.resolve(__dirname, 'src/assets'),
                '@pages': path.resolve(__dirname, 'src/pages'),
                "@package": path.resolve(__dirname, "package.json"),
            },
        },
    });
};
