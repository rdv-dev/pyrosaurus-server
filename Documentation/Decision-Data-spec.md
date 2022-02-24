# Decision Data Specification

Decisions are defined at a Dino's Species, so different Dinos may share the same Decision structure. Decisions trigger Movements based on what the Dino is able to sense.

The maximum size of the Decisions structure is 0x17D bytes. Each Decision structure is 0x14 bytes long. 
Therefore there is a maxiumum of 19 Decisions available.

## Decision Data Top-Level
Field|Size (bytes)|Description
---|---|---
Number of Decisions|1|The total number of Decision structures
Array of Decision Structures|19|See Decision Structure Table

## Decision Structure
Field|Size (bytes)|Description
---|---|---
Movement number|1|An identifier referring to the Movement to execute if this Decision evaluates true. See Movement Number Table for more information
Condition array|18|Array of bytes selecting Options for each Condition of the Decision. 

## Movement Number Table
The Movement number identifies which Movement to execute when a decision is evaluated to true. 
The Movement details are not included in the Team Entry Files, but they do work in the Test Arena, so the Game must have some implementation of each of these Movements. 
More research required here!

All User Created Movements start at ID number 7 and map to the Move Data Stucture in order starting at index 0.

## Internal Movement Descriptions
Reverse Engineering will reveal how exactly each of the internal Movements are implemented, but this has not yet been found.
In the meantime, we can use the descriptions from the manual.

Movement ID|Movement Name
---|---
0|Call
1|Don't Move
2|Wander
3|Move away
4|Move closer
5|Move North
6|Move South

### Call
A dino can roar to let other team members know where he is. They will come to his aid if they are trained to do so. This decision can take place at the same time asanother decision.  The dino may not roar immediately after he makes the decision to do so but will soon after. Because roars are so loud, dinos can hear roars twice as far as their hearing range.  Be careful, an opponent can also hear your roars.  A dino cannot shoot the same time he is roaring so his shots may be delayed if he is told to call.

### Don't Move
Stays where he is. He won't take any steps but he will turn to face his nearest opponent.

### Wander
Randomly wanders inside an area that has a diameter eight times the diameter of the zone of the last goal. A Wander decision can only be made when the dino has reached a goal unless it is the only decision. You can change the size of a zone in Movement Training.

### Move away
Moves directly away from the other dino.

### Move closer
Moves directly toward the other dino.

### Move North
Moves directly north toward enemy territory.

### Move South
Moves directly south away from enemy territory.

Each of the Decisions as presented in the Game is represented in the data as one byte.

## Condition Requirements
Each Condition is represented by one byte. Each Condition has one or two Options Lists available to it.
Each Option in the Option List is numbered where the first Option is 0, the second Option is 1, and so on.

Consider the following Condition:
```
He  <has>            <no>         legs.
    <does not have>  <two>
                     <four>
                     <any or no>
```

 * Selecting &lt;has&gt; and &lt;four&gt; will result in byte 0x20. 
 * Selecting &lt;does not have&gt; and &lt;two&gt; will result in byte 0x11.
 * Selecting &lt;has&gt; and &lt;any or no&gt; (this will turn off the Condition) will result in byte 0x30.


The second Options List is mapped to the high 4 bits and first Options list is mapped to the lower 4 bits of the byte. 
If a Condition only has one Option list, then that is the first Options list, and is mapped to the lower 4 bits of the byte and the higher 4 bits will stay 0.



The following is directly from the Pyrosaurus Manual.

Note: In the following descriptions, any words enclosed in < >
are options. One option must be selected from each column.  The
other options will be ignored.

For example, the first condition listed below could be read as
any of the following depending on which condition you select.

                  He is any friend.
                  He is any enemy.
                  He is my Queen.
                  He is enemy Queen.
                  I don't care who he is.



1.                He is  <any friend>.
                         <any enemy>.
                         <my Queen>.
                         <enemy Queen>.
                         <don't care>.

    Select the description that best fits the kind of dino that
    you want this decision to be concerned with.  All other
    dinos will not be considered.  If <don't care> is selected,
    then he will pick the nearest dino or ignore all dinos
    depending on the other conditions.



2.            He  <has>            <no>         legs.
                  <does not have>  <two>
                                   <four>
                                   <any or no>

    Narrow the decision process down to the type of dino.

    If "He is any friend" is selected for the first condition
    then the above condition is replaced with a selection that
    allows you to choose the friendly species. If you want to
    specify the species, click this selection and you will go to
    the Species Selection Screen where you can choose the
    species. If the species type doesn't matter, then press the
    CANCEL Button at the Species Selection screen and the
    species will be listed as "unknown".



3.       He is  <much smaller>  to  <much smaller>  than you.
                <smaller>           <smaller>
                <same size>         <same size>
                <larger>            <larger>
                <much larger>       <much larger>
                <any size>.

    Select his size relative to your dino's size.



4.          He  <is>      within  <fire range>.
                <is not>          <sight>.
                                  <hearing range>.
                                  <smelling range>.
                                  <any range>.

    How far away is he?  Dinos will not be noticed if they are
    past the farthest range.  If this condition is set to "is
    not within hearing range" and smelling range is the farthest
    range then dinos will be noticed if farther than the hearing
    range but within smelling range.

    A calling dino can be heard at twice the hearing range.  A
    dino that is not moving can only be heard at half the
    hearing range.

    If this condition is set to "any range" then this dino will
    notice anyone within his longest range.

    The selection "is not within any range" is not valid.

    A dino must be actually seen by your dino to be within
    sight. If your dino is turned away from the other dino, the
    other dino will not be seen no matter how close he is.

    Fire range is the range that your dino can shoot.

    Your dino will take most notice of the closest dino within
    the specified range.



5.        His skin is  <much thinner>  to  <much thinner>  than yours.
                       <thinner>           <thinner>
                       <the same>          <the same>
                       <thicker>           <thicker>
                       <much thicker>      <much thicker>
                       <don't care>

    The thicker a dino's skin, the more shots he can absorb. You
    may want to avoid any dinos that have thicker skin than
    yours and attack dinos with thinner skin. Or, you may want
    to gang up on a thick skinned opponent...




6.        My skin  <is>               <very thin>.
                   <is not>           <thin>.
                   <is thinner than>  <medium>.
                   <is thicker than>  <thick>.
                                      <very thick>.
                                      <don't care>.


    How thick is your dino's skin at the moment?  A dino's skin
    becomes thinner as he takes hits. When it becomes too thin,
    he dies.



7.        I am  <very tired>   to  <very tired>.
                <tired>            <tired>.
                <rested>           <rested>.
                <very rested>      <very rested>.
                                   <in any condition>.

    What kind of condition is your dino in?  When a dino runs
    and fights, he grows tired. As he becomes tired he won't be
    able to run and shoot as often.  The size of the dino's
    heart determines his endurance.  You should have your dino
    stop and rest if he gets too tired.




8.        My Queen  <is>      within  <fire range>  of him.
                    <is not>          <sight>
                                      <hearing range>
                                      <smelling range>
                                      <don't care>.

   Does the enemy know the location of your Queen?

   The Queen must be within one of your dino's senses before he
   is aware of her. If there is more than one Queen within
   range, only the nearest Queen will be noticed by this dino.



9.        My Queen  <is>      within  <sight>          of me.
                    <is not>          <hearing range>
                                      <smelling range>
                                      <don't care>.

   You might want to keep this dino close to the Queen to
   protect her from unwanted advances.

   This condition is not available if the first condition is
   set to "She is my Queen".




10.     His Queen  <is>      within  <fire range>  of me.
                   <is not>          <sight>
                                     <hearing range>
                                     <smelling range>
                                     <don't care>.


    Once this dino is aware of the enemy Queen, it may be time
    to go in for the kill or call for help.

    This condition is not available if the first condition is
    set to "She is enemy Queen".



11.       His full speed is  <less than>     to  <less than>     mine.
                             <same as>           <same as>
                             <greater than>      <greater than>
                             <don't care>.

    You may want to change your strategy if the enemy is faster
    than you. What if he is slower than you are? You could run
    away and look for his Queen.



12.             He  <is>      <attacking me>.
                    <is not>  <moving toward me>.
                              <moving perpendicular to me>.
                              <moving away from me>.
                              <attacking a friend>.
                              <doing anything>.


    You can react to another dino's movements here.




13.              He <is>                 calling.
                    <is not>
                    <may or may not be>

    If another team member is calling for help, do you want to
    come to his aid? What if an opponent is calling?



14.                Time <is>      <early>.
                        <is not>  <mid way>.
                                  <late>.
                                  <don't care>.

    This is contest time. You may want your team to act
    differently in the early part of a contest than near the
    end.



15.                Priority is  <very low>.
                                <low>.
                                <medium>.
                                <high>.
                                <very high>.

    How important is this decision?  There may be times when the
    dino could make more than one decision. When this happens,
    he will choose the one with the highest priority.

    If the movement for this decision is "Call" then Priority
    controls the frequency of his calls. If Priority is "very
    high" then he will call on the average of once every two
    footsteps.  If set to "very low" his calls will average once
    every 16 footsteps.  The other settings range between these
    two extremes.




The following selections are not part of the decision making
process but go into effect once this decision is made:



16.               Food is <not important>.
                          <important>.
                          <very important>.


    Dinos gain skin thickness and become less tired when they
    eat food.  If food is not important he will only eat if
    it is very close and he is very weak.  If food is very
    important, he may take a long detour to get to it.

    The disadvantage of going to food is that it will delay
    progress which could be crucial when going to the aid of
    someone in trouble.

    Food is not available in all arenas.  It is a good
    idea to assume that there is always food so that your
    training will work correctly when food is available.



17.                   I should  <not move>.
                                <creep>.
                                <walk>.
                                <run>.

    This is the speed that the dino should travel if he makes this
    decision. If the dino is too tired to run, he will walk.

    If this decision is set to Pack (see below) then the dino
    will use whatever speed is needed to stay with the pack.


18.                       <Pack>.
                          <Don't Pack>.


     Only available if the movement type is Other-Mobile. If
     set to Pack then the dino will try to stay as close as he
     can to his goal.  If he is slower then the rest of the
     pack, then they will slow down to let him catch up.

     If set to Don't Pack then he will go to his leader's next
     goal without trying to stay with the leader. If he is
     faster than the rest of the pack, he will move in front of
     the pack. If he is slower than the pack, they will not slow
     down.  This is a good way to send scout dinos out ahead of
     the pack.

     A dino that has selected a Don't Pack decision will not
     follow the conventions used when packing. This means that
     he may get in his leader's way or not move out of the way
     if he is sitting on his leader's goal.




