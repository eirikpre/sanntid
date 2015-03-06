with Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;
use  Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;

procedure exercise7 is

    Count_Failed    : exception;    -- Exception to be raised when counting fails
    Gen             : Generator;    -- Random number generator

    protected type Transaction_Manager (N : Positive) is
        entry Finished;
        procedure Signal_Abort;
        entry Wait_Until_Aborted;
    private
        Finished_Gate_Open  : Boolean := False;
        Aborted             : Boolean := False;
        
    end Transaction_Manager;

    protected body Transaction_Manager is
        entry Finished when Finished_Gate_Open or Finished'Count = N is
        begin
            Finished_Gate_Open := True;
            if Finished'Count = 0 then
                Finished_Gate_Open := False;
            end if;                       
        end Finished;

        entry Wait_Until_Aborted when Aborted is
        begin
                if Wait_Until_Aborted'Count = 0 then
                    Aborted := False;
                end if;
        end Wait_Until_Aborted;

        procedure Signal_Abort is
        begin
            Aborted := True;
        end Signal_Abort;

    end Transaction_Manager;



    
    function Unreliable_Slow_Add (x : Integer) return Integer is
    Error_Rate : Constant := 0.25;  -- (between 0 and 1)
    begin
        if  Random(Gen) < Error_rate then
            delay Duration(1);
            raise Count_Failed;
        else 
            delay Duration(2);
            return (x+10);
        end if;
    end Unreliable_Slow_Add;


    task type Transaction_Worker (Initial : Integer; Manager : access Transaction_Manager);
    task body Transaction_Worker is
        Num         : Integer   := Initial;
        Prev        : Integer   := Num;
        Round_Num   : Integer   := 0;
    begin
        Put_Line ("Worker" & Integer'Image(Initial) & " started");

        loop
            Put_Line ("Worker" & Integer'Image(Initial) & " started round" & Integer'Image(Round_Num));
            Round_Num := Round_Num + 1;





            

            select
                Manager.Wait_Until_Aborted;
                Num := Prev + 5;
                Put_Line ("  Worker" & Integer'Image(Initial) &
                             " forward ER from " & Integer'Image(Prev) &
                             " to" & Integer'Image(Num));

                

            then abort
                begin
                Num := Unreliable_Slow_Add(Num);
                
                exception 
                    when Count_Failed =>
                        Manager.Signal_Abort;
                end;
                
                  Put_Line ("  Worker" & Integer'Image(Initial) & " comitting" & Integer'Image(Num)); 
                Manager.Finished;
            end select;

            

            delay 0.5;         
            Prev := Num;
        end loop;
    end Transaction_Worker;


    Manager : aliased Transaction_Manager (3);
    Worker_1 : Transaction_Worker (0, Manager'Access);
    Worker_2 : Transaction_Worker (1, Manager'Access);
    Worker_3 : Transaction_Worker (2, Manager'Access);

begin
    
    Reset(Gen); -- Seed the random number generator
end exercise7;



