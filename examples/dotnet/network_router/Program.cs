using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet network_router Group",
    });

    var network = new Network("example-network", new NetworkArgs
    {
        Name = "Example DotNet network_router Network",
    });

    var res = new NetworkRouter("test-network-router", new NetworkRouterArgs
    {
        NetworkId = network.Id,
        PeerGroups = { group.Id },
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.NetworkId,
    };
});
