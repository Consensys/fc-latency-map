import React, { useState, useEffect } from "react";
import Head from "next/head";
import dynamic from "next/dynamic";
import getConfig from "next/config";

import styles from "@src/styles/Global.module.css";
import Header from "@src/components/Header";

const { publicRuntimeConfig, serverRuntimeConfig } = getConfig();

interface Props {
  data: any;
}

const Home = (props: Props) => {
  const { data } = props;
  const MapWithNoSSR = dynamic(() => import("@src/components/map/Map"), {
    ssr: false,
  });

  const [size, setSize] = useState({
    x: 0,
    y: 0,
  });

  const updateSize = () => {
    if (size.x != window.innerWidth || size.y != window.innerHeight) {
      setSize({
        x: window.innerWidth,
        y: window.innerHeight,
      });
    }
  };

  useEffect(() => {
    setSize({
      x: window.innerWidth,
      y: window.innerHeight,
    });
  }, []);

  useEffect(() => (window.onresize = updateSize), []);

  return (
    <div className={styles.container}>
      <Head>
        <title>Filecoin Latency Map</title>
        <meta name="description" content={publicRuntimeConfig.app.name} />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Header />
      <div
        id="map"
        style={{
          width: "100%",
          minHeight: "600px",
        }}
      >
        <MapWithNoSSR data={data} height={size.x} width={size.y} />
      </div>
    </div>
  );
};

export async function getServerSideProps() {
  // const data = await import(
  //   serverRuntimeConfig.path.exportsMeasures + "export.json"
  // );
  const data = await import("../../data/data01.json");

  return {
    props: {
      data: JSON.stringify(data),
    },
  };
}

export default Home;
