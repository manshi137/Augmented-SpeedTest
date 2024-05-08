# read from ping_reply_map.csv
import csv
import numpy as np
import scipy.stats as stats
from datetime import datetime
import os
import pandas as pd

hops = 3
alpha = 0.05

def calculate_time_difference(request_time, reply_time):
    request_time = datetime.strptime(request_time, "%H:%M:%S.%f")
    reply_time = datetime.strptime(reply_time, "%H:%M:%S.%f")
    time_diff = reply_time - request_time
    return time_diff.total_seconds()

def run_ttest():
    # read from ping_reply.csv
    ping_reply_map = []
    download = [[] for i in range(hops+1)]
    upload = [[] for i in range(hops+1)]
    idle = [[] for i in range(hops+1)]

    congested_d = None
    congested_u = None

    



    

    # Read the CSV file into a DataFrame
    df = pd.read_csv('ping_reply.csv')

    # Convert the RequestTime column to datetime format
    df['RequestTime'] = pd.to_datetime(df['RequestTime'], format='%H:%M:%S.%f')

    # Sort the DataFrame by the RequestTime column
    df_sorted = df.sort_values(by='RequestTime')

    # Convert the RequestTime column back to the desired format
    df_sorted['RequestTime'] = df_sorted['RequestTime'].dt.strftime('%H:%M:%S.%f')

    # Write the sorted DataFrame to a new CSV file
    # df_sorted.to_csv('sorted_ping_reply.csv', index=False)

    df_sorted.to_csv('sorted_ping_reply.csv', index=False)

    with open('sorted_ping_reply.csv', 'r') as file:
        reader = csv.DictReader(file)
        for row in reader:
            request_time = row['RequestTime']
            reply_time = row['ReplyTime']
            time_diff = calculate_time_difference(request_time, reply_time)
            # print(f"Time difference: {time_diff} seconds")
            if(row['Download/Upload/Idle'] == 'download'):
                download[int(row['TTL'])].append(time_diff)
            elif(row['Download/Upload/Idle'] == 'upload'):
                upload[int(row['TTL'])].append(time_diff)
            elif(row['Download/Upload/Idle'] == 'idle'):
                idle[int(row['TTL'])].append(time_diff)
            else:
                print("onooo")

    # -----------------checkkk length ----------------
    print(download[1])
    print(download[2])
    print(download[3])
    n = min(len(download[1]),len(download[2]),len(download[3]), len(upload[1]), len(upload[2]), len(upload[3]))
    n = min(n, len(idle[1]), len(idle[2]), len(idle[3]))
    print("n ==", n)
    # truncate the list to the same length
    for h in range(1, hops+1):
        download[h] = download[h][:n]
        upload[h] = upload[h][:n]
        idle[h] = idle[h][:n]
        # print(f"download[{h}]: {len(download[h])}, upload[{h}]: {len(upload[h])}, idle[{h}]: {len(idle[h])}")

    # convert the cumulative time to time difference
    print("cumulative download")
    print(download[1])
    print(download[2])
    print(download[3])

    for hop in range(hops, 1, -1):
        download[hop] = [download[hop][i] - download[hop-1][i] for i in range(n)]
        upload[hop] = [upload[hop][i] - upload[hop-1][i] for i in range(n)]
        idle[hop] = [idle[hop][i] - idle[hop-1][i] for i in range(n)]

    print("non cumulative download")
    print(download[1])
    print(download[2])
    print(download[3])
    print("\n---------------------------------------------------------------")
    p_values_map_d = {}
    t_statistics_map_d = {}
    p_values_map_u = {}
    t_statistics_map_u = {}
    # download vs idle
    for hop in range(1, hops+1):
        t_statistic, p_value = stats.ttest_ind(download[hop], idle[hop])
        t_statistic, p_value = stats.ttest_ind(download[hop], idle[hop])
        p_values_map_d[hop] = p_value
        t_statistics_map_d[hop] = t_statistic
        print("hop= ", hop, "download vs idle")
        print("t_statistic= ", t_statistic)
        print("p_value= ", p_value)

    print("------------------------------------------- \n")
    # upload vs idle
    for hop in range(1, hops+1):
        t_statistic, p_value = stats.ttest_ind(upload[hop], idle[hop])
        p_values_map_u[hop] = p_value
        t_statistics_map_u[hop] = t_statistic
        print("hop= ", hop, "upload vs idle")
        print("t_statistic= ", t_statistic)
        print("p_value= ", p_value)
        
    print("\n---------------------------------------------------------------")

    max_tstat_d = float('-inf')
    max_tstat_hop_d = None
    congested_d = []

    for hop in range(1, hops+1):
        if hop in p_values_map_d and p_values_map_d[hop] < alpha:
            if hop in t_statistics_map_d and t_statistics_map_d[hop] > max_tstat_d:
                congested_d = [hop]
                max_tstat_d = t_statistics_map_d[hop]
                max_tstat_hop_d = hop
            elif hop in t_statistics_map_d and t_statistics_map_d[hop] == max_tstat_d:
                congested_d.append(hop)

        max_tstat_u = float('-inf')
        max_tstat_hop_u = None
        congested_u = []

    for hop in range(1, hops+1):
        if hop in p_values_map_u and p_values_map_u[hop] < alpha:
            if hop in t_statistics_map_u and t_statistics_map_u[hop] > max_tstat_u:
                congested_u = [hop]
                max_tstat_u = t_statistics_map_u[hop]
                max_tstat_hop_u = hop
            elif hop in t_statistics_map_u and t_statistics_map_u[hop] == max_tstat_u:
                congested_u.append(hop)

    print("Congested hops for download:", congested_d)
    print("Hop with highest t_statistic value for download:", max_tstat_hop_d)
    print("Corresponding t_statistic value for download:", max_tstat_d)

    print("Congested hops for upload:", congested_u)
    print("Hop with highest t_statistic value for upload:", max_tstat_hop_u)
    print("Corresponding t_statistic value for upload:", max_tstat_u)

    print("\n")
    print("-------------------------------------------------------------test complete-----------------------------------------------------------")




run_ttest()

