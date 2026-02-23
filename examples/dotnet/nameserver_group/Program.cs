using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet nameserver_group Group",
    });

    var res = new NameserverGroup("test-nameserver-group", new NameserverGroupArgs
    {
        Name = "Pulumi DotNet NS Group",
        Description = "Pulumi DotNet NS Group",
        Enabled = true,
        Nameservers = { new NameserverGroupNameserverArgs { Ip = "1.1.1.1", Port = 53 } },
        Groups = { group.Id },
        Primary = true,
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.Name,
    };
});
