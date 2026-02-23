using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet network_resource Group",
    });

    var network = new Network("example-network", new NetworkArgs
    {
        Name = "Example DotNet network_resource Network",
    });

    var res = new NetworkResource("test-network-resource", new NetworkResourceArgs
    {
        Name = "Pulumi DotNet Net Res",
        Address = "10.20.0.0/24",
        NetworkId = network.Id,
        Groups = { group.Id },
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Name,
    };
});
