import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Test and Run Strategies',
    Svg: require('@site/static/img/logo.svg').default,
    description: (
      <>
        TODO
      </>
    ),
  },
  {
    title: 'Connected to the world',
    Svg: require('@site/static/img/logo.svg').default,
    description: (
      <>
        TODO
      </>
    ),
  },
  {
    title: 'Scalable and Production ready',
    Svg: require('@site/static/img/logo.svg').default,
    description: (
      <>
        TODO
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
