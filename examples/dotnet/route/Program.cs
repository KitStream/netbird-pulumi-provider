using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet route Group",
    });

    var res = new Route("test-route", new RouteArgs
    {
        Description = "Pulumi DotNet Route",
        Enabled = true,
        Network = "10.0.0.0/24",
        NetworkId = "test-route",
        PeerGroups = { group.Id },
        Groups = { group.Id },
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Description,
    };
});
