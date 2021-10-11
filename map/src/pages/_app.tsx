import Head from "next/head";
import "antd/dist/antd.css";
import "@src/styles/globals.css";
import type { AppProps } from "next/app";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <div>
      <Head>
        {/* <link
          rel="apple-touch-icon"
          sizes="196x196"
          href="./favicon-196x196.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="196x196"
          href="./favicon-196x196.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="32x32"
          href="./favicon-32x32.png"
        />
        <link
          rel="icon"
          type="image/png"
          sizes="16x16"
          href="./favicon-16x16.png"
        /> */}
      </Head>
      <Component {...pageProps} />
    </div>
  );
}
export default MyApp;
