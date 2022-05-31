# co2-exporter

A prometheus exporter collecting the actual CO2 emission from France's electricity consumption. It is based on the [RTE eco2mix app](https://www.rte-france.com/eco2mix). 

# Build

```
docker build -t co2-exporter .
```

# Deploy

Configuration du deployment dans le `deployment.yaml`.  
```
kubectl apply -f deployment.yaml
```

# prometheus-adapter

Configuration dans `rules.yaml`. 
```
helm upgrade --install prometheus-adapter prometheus-community/prometheus-adapter -n monitoring-system -f rules.yml
```

# horizontal pod autoscaler
Configuration dans `hpa.yaml`
```
kubectl apply -f hpa.yaml
```

# Liens utiles

- [RTE eco2mix](https://www.rte-france.com/eco2mix)
- [Github prometheus-adapter](https://github.com/kubernetes-sigs/prometheus-adapter)
- 