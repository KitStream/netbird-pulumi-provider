using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet network Group",
    });

    var res = new Network("test-network", new NetworkArgs
    {
        Name = "Pulumi DotNet Network",
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Name,
    };
});
