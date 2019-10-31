import matplotlib
matplotlib.use('tkagg')
import matplotlib.pyplot as plt
import re



ns = [x for x in range(0, 501, 50)]
data = {}
with open('results_parsed') as f:
    for line in f:
        m = re.split(' ', line)
        
        l = data.get(m[0], {})
        size = int(m[1])
        time = int(m[2])
        l.update({size:time})
        
        data.update({m[0]: l})

mydata = {}
for label in data:
    points = data[label]
    nums = [points[key] for key in sorted(points.keys())]
    mydata.update({label: nums})

for (label, values) in mydata.items():
    plt.plot(ns, values, label=label)
plt.legend()
plt.show()