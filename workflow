Connect to influxdb db
Read Configuration file ( Regions, types )
Generate quires as needed
Verify at least 1 data point from each back data day for each query

-> Prediction Group
-> Wait until Backfilled data
-> Begin Prediction Process (write results to db)

AWS Group
-> Set up AWS endpoint
-> Pass Initial Query (if needed)
-> Store last checked time
-> Sleep for Period



Prediction Requirements (0 means not possible)
Initial:
  Below X Price Y Percent of Time
Advanced:
  Below X Price Y Percent of Time
  Spikes mitigated by spike duration total cost < X by hour (??)

Write Cheapest solution for all queries  
