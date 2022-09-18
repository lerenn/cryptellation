import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Connected to the world',
    Svg: require('@site/static/img/cloud.svg').default,
    description: (
      <>
        The Cryptellation system aims to connect your strategies to different
        marketplaces and have different newsfeed to react to any opportunity.
      </>
    ),
  },
  {
    title: 'Test and run strategies',
    Svg: require('@site/static/img/repeat-candlesticks.svg').default,
    description: (
      <>
        With <b>backtesting</b>, <b>forward testing</b> and <b>live running</b> capabilities, you can test
        your algorithms on any market conditions to control how it can react and
        measure its <b>profitability</b>.
      </>
    ),
  },
  {
    title: 'Scalable and production ready',
    Svg: require('@site/static/img/production.svg').default,
    description: (
      <>
        With <b>unit</b>, <b>integration</b> and <b>end-to-end</b> tests, the system aims
        for <b>reliability</b> and <b>security</b>. Also, you can plug the system
        into your favorite administrator tools (Kubernetes, Prometheus, and more).
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
