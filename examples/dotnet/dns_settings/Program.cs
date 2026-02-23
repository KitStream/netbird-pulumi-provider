using System.Collections.Generic;
using Pulumi;
using KitStream.Pulumi.Netbird;

return await Deployment.RunAsync(() => 
{
    var group = new Group("example-group", new GroupArgs
    {
        Name = "Example DotNet dns_settings Group",
    });

    var res = new DnsSettings("test-dns-settings", new DnsSettingsArgs
    {
        DisabledManagementGroups = { group.Id },
    });

    return new Dictionary<string, object?>
    {
        ["resourceName"] = res.DisabledManagementGroups,
    };
});
